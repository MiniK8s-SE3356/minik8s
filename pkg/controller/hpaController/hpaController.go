package hpacontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/hpa"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

type HPAController struct {
}

func NewHPAController() *(HPAController) {
	fmt.Printf("New HPAController...\n")
	return &HPAController{}
}

func checkMatchedPod(podLabels map[string]string, selector map[string]string) bool {
	for k, v := range podLabels {
		if selector[k] == v {
			return true
		}
	}

	return false
}

func getPodsFromServer() ([]pod.Pod, error) {
	var result []pod.Pod
	var podMap map[string]pod.Pod
	jsonData, err := httpRequest.GetRequest(url.RootURL + url.GetAllPod)
	if err != nil {
		fmt.Println("error in get request")
		return result, err
	}

	fmt.Println(string(jsonData))
	err = json.Unmarshal([]byte(jsonData), &podMap)
	if err != nil {
		fmt.Println("failed to unmarshal")
		fmt.Println(jsonData)
		return result, err
	}

	for _, v := range podMap {
		result = append(result, v)
	}

	return result, nil
}

func getHPAsFromServer() ([]hpa.HPA, error) {
	var result []hpa.HPA

	jsonData, err := httpRequest.GetRequest(url.RootURL + url.GetHPA)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("error in get request")
		return result, err
	}

	var rsMap map[string]hpa.HPA
	err = json.Unmarshal([]byte(jsonData), &rsMap)
	if err != nil {
		fmt.Println("failed to unmarshal")
		fmt.Println(jsonData)
		return result, err
	}

	for _, hpa := range rsMap {
		result = append(result, hpa)
	}

	return result, nil
}

func applyPod(pod yaml.PodDesc) error {
	fmt.Println("apply pod", pod)

	req := make(map[string]interface{})
	req["namespace"] = process.DefaultNamespace
	req["podDesc"] = pod

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("failed to translate into json")
		return err
	}
	// fmt.Println(podDesc.Spec.Containers)
	result, err := httpRequest.PostRequest(url.RootURL+url.AddPod, jsonData)
	if err != nil {
		fmt.Println("error when post request")
		return err
	}

	fmt.Println(result)

	return nil
}

func removePod(namespace string, name string) error {
	fmt.Println("remove", namespace, name)
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return err
	}
	result, err := httpRequest.PostRequest(url.RootURL+url.RemovePod, jsonData)
	if err != nil {
		fmt.Println("error in delete pod ", err.Error())
		return err
	}
	fmt.Println(result)

	return nil
}

func updateHPA(hpa hpa.HPA) error {
	req := make(map[string]interface{})
	req["namespace"] = process.DefaultNamespace
	req["hpa"] = hpa

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("failed to translate into json")
		return err
	}
	result, err := httpRequest.PostRequest(url.RootURL+url.UpdateHPA, jsonData)
	if err != nil {
		fmt.Println("error when post request")
		return err
	}

	fmt.Println(result)

	return nil
}

func rscTask() {
	// 从APIserver获取全体pod和全体hpa
	pods, err := getPodsFromServer()
	if err != nil {
		fmt.Println("failed to get pod from server", err)
		return
	}

	hpas, err := getHPAsFromServer()
	if err != nil {
		fmt.Println("failed to get hpa from server", err)
	}

	hpaNames := make(map[string]bool, 0)
	// 利用label找到hpa对应的pod
	for _, hpa := range hpas {
		hpaNames[hpa.Metadata.Name] = true
		var matchedPod []pod.Pod

		// 先检查hpa是否到了timeInterval，如果没有就跳过
		if time.Since(hpa.Status.LastUpdateTime) < hpa.Spec.TimeInterval {
			continue
		}
		hpa.Status.LastUpdateTime = time.Now()

		for _, p := range pods {
			if p.Metadata.Labels["hpa"] == hpa.Metadata.Name {
				// if checkMatchedPod(p.Metadata.Labels, hpa.Spec.Selector.MatchLabels) {
				matchedPod = append(matchedPod, p)
			}
		}

		// 比对pod数量和期望数量
		change := 0

		if len(matchedPod) > hpa.Spec.MaxReplicas {
			// 大于maxReplicas
			// delete some pod
			change = -1
		} else if len(matchedPod) < hpa.Spec.MinReplicas {
			// 小于minReplicas
			// add some pod
			// for i := 0; i < hpa.Spec.MinReplicas-len(matchedPod); i++ {
			change = 1

		} else {
			// 根据pod平均资源使用率计算
			averageCPU := 0.0
			averageMem := 0.0
			for _, p := range matchedPod {
				averageCPU += p.Status.CPUUsage
				averageMem += p.Status.MemoryUsage
			}
			averageCPU /= float64(len(matchedPod))
			averageMem /= float64(len(matchedPod))

			if (averageCPU > hpa.Spec.Metrics.CPUPercent || averageMem > hpa.Spec.Metrics.MemPercent) && len(matchedPod)+1 <= hpa.Spec.MaxReplicas {
				change = 1
			} else if (averageCPU < hpa.Spec.Metrics.CPUPercent || averageMem < hpa.Spec.Metrics.MemPercent) && len(matchedPod)-1 >= hpa.Spec.MinReplicas {
				change = -1
			}
		}

		if change == 1 {
			var podDesc yaml.PodDesc
			podDesc.ApiVersion = "v1"
			podDesc.Kind = "Pod"
			newId, _ := idgenerate.GenerateID()
			podDesc.Metadata.Name = hpa.Metadata.Name + "-" + newId[:8]
			podDesc.Metadata.Labels = make(map[string]string)
			podDesc.Metadata.Labels["hpa"] = hpa.Metadata.Name
			podDesc.Spec.Containers = hpa.Spec.Template.Spec.Containers
			err := applyPod(podDesc)
			if err != nil {
				fmt.Println("failed to apply pod")
			}
		} else if change == -1 {
			for i := 0; i < 1; i++ {
				ns := matchedPod[i].Metadata.Namespace
				name := matchedPod[i].Metadata.Name
				err := removePod(ns, name)
				if err != nil {
					fmt.Println("failed to delete pod", err)
				}
			}
		}
	}

	// 处理孤儿pod
	// 从hpa创建的Pod都带有与hpa对应的label
	for _, p := range pods {
		value1, ok1 := p.Metadata.Labels["hpa"]
		if ok1 {
			_, ok2 := hpaNames[value1]
			if !ok2 {
				removePod(p.Metadata.Namespace, p.Metadata.Name)

			}
		}
	}
}

func (sc *HPAController) Init() {
	fmt.Printf("Init HPAController ...\n")

}

func (sc *HPAController) Run() {
	fmt.Printf("Run HPAController ...\n")

	for {

		rscTask()
		<-time.After(20 * time.Second)
	}
}
