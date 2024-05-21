package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddPod(namespace string, desc *yaml.PodDesc) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	// 构建然后转json
	pod := &pod.Pod{}
	pod.APIVersion = desc.ApiVersion
	pod.Kind = desc.Kind
	pod.Metadata.Name = desc.Metadata.Name
	pod.Metadata.Namespace = namespace
	pod.Metadata.UUID = id
	pod.Metadata.Labels = desc.Metadata.Labels
	pod.Spec = desc.Spec
	// for _, c := range desc.Spec.Containers {
	// var tmp container.Container
	// pod.Spec.Containers = append(pod.Spec.Containers, )
	// }

	value, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 先查看一下key是否已经存在
	tmp, err := EtcdCli.Exist(podPrefix + namespace + "/" + pod.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check existence in etcd")
		return "failed to check existence in etcd", err
	}
	if tmp {
		fmt.Println("pod has existed")
		return "pod has existed", nil
	}
	// 然后存入etcd
	err = EtcdCli.Put(podPrefix+namespace+"/"+pod.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	msgBody := make(map[string]interface{})
	msgBody["type"] = "create_pod"
	msgBody["content"] = pod
	jsonData, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("failed to construct msgBody")
		return "failed to construct msgBody", err
	}
	Mq.Publish("minik8s", "scheduler", "application/json", jsonData)

	return "add pod to minik8s", nil
}

func RemovePod(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	existed, err := EtcdCli.Exist(podPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "target not exist", nil
	}

	err = EtcdCli.Del(podPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del pod success", nil
}

func UpdatePod(namespace string, pod pod.Pod) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	name := pod.Metadata.Name

	existed, err := EtcdCli.Exist(podPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "target not exist", nil
	}

	value, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}
	err = EtcdCli.Put(podPrefix+namespace+"/"+name, string(value))
	if err != nil {
		fmt.Println("failed to put into etcd")
		return "failed to put into etcd", err
	}

	return "update success", nil
}

func GetPod(namespace string, name string) (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	var r pod.Pod
	result := make(map[string]interface{}, 0)

	existed, err := EtcdCli.Exist(podPrefix + namespace + "/" + name)
	if err != nil {
		return result, err
	}
	if !existed {
		return result, nil
	}

	tmp, err := EtcdCli.Get(podPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, nil
	}

	err = json.Unmarshal(tmp, &r)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return result, nil
	}

	result[podPrefix+namespace+"/"+name] = r

	return result, nil
}

func GetPods(namespace string) (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]interface{}, 0)

	pairs, err := EtcdCli.GetWithPrefix(podPrefix + namespace)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			result[p.Key] = tmp
		}
	}

	return result, nil
}

func GetAllPods() (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]interface{}, 0)

	pairs, err := EtcdCli.GetWithPrefix(podPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			result[p.Key] = tmp
		}
	}

	return result, nil
}

// func removePodByNode(nodeIP string) (string, error) {
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	pairs, err := EtcdCli.GetWithPrefix(podPrefix)
// 	if err != nil {
// 		fmt.Println("failed to get from etcd")
// 		return "", err
// 	}

// 	var deleteList []string
// 	for _, p := range pairs {
// 		var pod pod.Pod
// 		err := json.Unmarshal([]byte(p.Value), &pod)
// 		if err != nil {
// 			fmt.Println("failed to translate into json")
// 			continue
// 		} else {
// 			pod.Spec
// 		}
// 	}
// }

func DescribePod(pod string, name string) (string, error) {
	return "", nil
}

func DescribePods(pod string) (string, error) {
	return "", nil
}

func DescribeAllPods() (string, error) {
	return "", nil
}
