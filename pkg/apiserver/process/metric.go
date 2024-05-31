package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
)

func GetMetricPoint() (map[string]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]string)

	// 先检查pod
	pairs, err := EtcdCli.GetWithPrefix(podPrefix)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	for _, p := range pairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 这里也包括具体路径
		port, ok := tmp.Metadata.Labels["metric_port"]
		if ok {
			result[tmp.Metadata.Name] = tmp.Status.PodIP + ":" + port
		}
	}

	// 再看node
	pairs, err = EtcdCli.GetWithPrefix(nodePrefix)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	for _, p := range pairs {
		var tmp node.Node
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		port, ok := tmp.Metadata.Labels["metric_port"]
		if ok {
			result[tmp.Metadata.Name] = tmp.Status.Ip + ":" + port
		}
	}

	return result, nil
}
