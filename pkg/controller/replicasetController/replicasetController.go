package replicasetcontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/replicaset"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/config"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

type ReplicasetController struct {
}

func NewReplicasetController() *(ReplicasetController) {
	fmt.Printf("New ReplicasetController...\n")
	return &ReplicasetController{}
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
	jsonData, err := httpRequest.GetRequest(config.HTTPURL + url.GetAllPod)
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

func getReplicasetsFromServer() ([]replicaset.Replicaset, error) {
	var result []replicaset.Replicaset

	jsonData, err := httpRequest.GetRequest(config.HTTPURL + url.GetReplicaset)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("error in get request")
		return result, err
	}

	var rsMap map[string]replicaset.Replicaset
	err = json.Unmarshal([]byte(jsonData), &rsMap)
	if err != nil {
		fmt.Println("failed to unmarshal")
		fmt.Println(jsonData)
		return result, err
	}

	for _, rs := range rsMap {
		result = append(result, rs)
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
	result, err := httpRequest.PostRequest(config.HTTPURL+url.AddPod, jsonData)
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
	result, err := httpRequest.PostRequest(config.HTTPURL+url.RemovePod, jsonData)
	if err != nil {
		fmt.Println("error in delete pod ", err.Error())
		return err
	}
	fmt.Println(result)

	return nil
}

func rscTask() {
	// 从APIserver获取全体pod和全体replicaset
	pods, err := getPodsFromServer()
	if err != nil {
		fmt.Println("failed to get pod from server", err)
		return
	}

	replicasets, err := getReplicasetsFromServer()
	if err != nil {
		fmt.Println("failed to get replicaset from server", err)
	}

	replicasetNames := make(map[string]bool, 0)
	// 利用label找到replicaset对应的pod
	for _, rs := range replicasets {
		replicasetNames[rs.Metadata.Name] = true
		var matchedPod []pod.Pod

		for _, p := range pods {
			if p.Metadata.Labels["replicaset"] == rs.Metadata.Name {
				matchedPod = append(matchedPod, p)
			}
		}
		fmt.Println(pods, rs, matchedPod)

		// 比对pod数量和期望数量
		fmt.Println("...................[]................")
		rsJSON, _ := json.Marshal(rs)
		podsJSON, _ := json.Marshal(pods)
		matchedPodJSON, _ := json.Marshal(matchedPod)
		fmt.Println(string(rsJSON))
		fmt.Println(string(podsJSON))
		fmt.Println(string(matchedPodJSON))
		fmt.Println("...................[]................")
		if len(matchedPod) > rs.Spec.Replicas {
			// delete some pod
			for i := 0; i < len(matchedPod)-rs.Spec.Replicas; i++ {
				ns := matchedPod[i].Metadata.Namespace
				name := matchedPod[i].Metadata.Name
				err := removePod(ns, name)
				if err != nil {
					fmt.Println("failed to delete pod", err)
				}
			}
		} else if len(matchedPod) < rs.Spec.Replicas {
			// add some pod
			for i := 0; i < rs.Spec.Replicas-len(matchedPod); i++ {
				var podDesc yaml.PodDesc
				podDesc.ApiVersion = "v1"
				podDesc.Kind = "Pod"
				newId, _ := idgenerate.GenerateID()
				podDesc.Metadata.Name = rs.Metadata.Name + "-" + newId[:8]
				podDesc.Metadata.Labels = make(map[string]string)
				podDesc.Metadata.Labels["replicaset"] = rs.Metadata.Name
				podDesc.Spec.Containers = rs.Spec.Template.Spec.Containers
				err := applyPod(podDesc)
				if err != nil {
					fmt.Println("failed to apply pod")
				}
			}
		}
	}

	// 处理孤儿pod
	// 从replicaset创建的Pod都带有与replicaset对应的label
	for _, p := range pods {
		value1, ok1 := p.Metadata.Labels["replicaset"]
		if ok1 {
			_, ok2 := replicasetNames[value1]
			if !ok2 {
				removePod(p.Metadata.Namespace, p.Metadata.Name)

			}
		}
	}

}

func (sc *ReplicasetController) Init() {
	fmt.Printf("Init ReplicasetController ...\n")

}

func (sc *ReplicasetController) Run() {
	fmt.Printf("Run ReplicasetController ...\n")

	for {

		rscTask()
		<-time.After(15 * time.Second)
	}
}
