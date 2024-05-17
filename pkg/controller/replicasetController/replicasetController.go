package replicasetcontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/replicaset"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
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

	jsonData, err := httpRequest.GetRequest(url.RootURL + url.GetPod)
	if err != nil {
		fmt.Println("error in get request")
		return result, err
	}

	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		fmt.Println("failed to unmarshal")
		fmt.Println(jsonData)
		return result, err
	}

	return result, err
}

func getReplicasetsFromServer() ([]replicaset.Replicaset, error) {
	var result []replicaset.Replicaset

	jsonData, err := httpRequest.GetRequest(url.RootURL + url.GetReplicaset)
	if err != nil {
		fmt.Println("error in get request")
		return result, err
	}

	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		fmt.Println("failed to unmarshal")
		fmt.Println(jsonData)
		return result, err
	}

	return result, err
}

func rscTask() {
	// 从APIserver获取全体pod和全体replicaset
	pods, err := getPodsFromServer()
	if err != nil {
		fmt.Println("failed to get pod from server")
		return
	}

	replicasets, err := getReplicasetsFromServer()
	if err != nil {
		fmt.Println("failed to get replicaset from server")
	}

	// 利用label找到replicaset对应的pod
	for _, rs := range replicasets {
		var matchedPod []pod.Pod

		for _, pod := range pods {
			if checkMatchedPod(pod.Metadata.Labels, rs.Spec.Selector.MatchLabels) {
				matchedPod = append(matchedPod, pod)
			}
		}

		// 比对pod数量和期望数量

		if len(matchedPod) > rs.Spec.Replicas {
			// delete some pod
		} else if len(matchedPod) < rs.Spec.Replicas {
			// add some pod
		}
	}

	// 处理孤儿pod
}

func (sc *ReplicasetController) Init() {
	fmt.Printf("Init ReplicasetController ...\n")

}

func (sc *ReplicasetController) Run() {
	fmt.Printf("Run ReplicasetController ...\n")

	for {
		time.After(10 * time.Second)

		rscTask()
	}
}
