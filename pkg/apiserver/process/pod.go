package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
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
	pod_ := &pod.Pod{}
	pod_.APIVersion = desc.ApiVersion
	pod_.Kind = desc.Kind
	pod_.Status.Phase = pod.PodPending
	pod_.Metadata.Name = desc.Metadata.Name
	pod_.Metadata.Namespace = namespace
	pod_.Metadata.UUID = id
	pod_.Metadata.Labels = desc.Metadata.Labels
	pod_.Spec = desc.Spec

	//! Avoid the same container name between different pods
	for i, c := range pod_.Spec.Containers {
		pod_.Spec.Containers[i].Name = pod_.Metadata.UUID + "-" + c.Name
	}

	// for _, c := range desc.Spec.Containers {
	// var tmp container.Container
	// pod.Spec.Containers = append(pod.Spec.Containers, )
	// }

	value, err := json.Marshal(pod_)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 先查看一下key是否已经存在
	tmp, err := EtcdCli.Exist(podPrefix + namespace + "/" + pod_.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check existence in etcd")
		return "failed to check existence in etcd", err
	}
	if tmp {
		fmt.Println("pod has existed")
		return "pod has existed", nil
	}
	// 然后存入etcd
	err = EtcdCli.Put(podPrefix+namespace+"/"+pod_.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	msgBody := make(map[string]interface{})
	msgBody["type"] = "create_pod"
	msgBody["content"] = pod_
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
	var r pod.Pod

	existed, err := EtcdCli.Exist(podPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "target not exist", nil
	}

	tmp, err := EtcdCli.Get(podPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "failed to get from etcd", nil
	}

	err = json.Unmarshal(tmp, &r)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return "failed to unmarshal", nil
	}
	// 这里有一个假设是一定能够删除成功
	msgBody := make(map[string]interface{})
	msgBody["type"] = message.RemovePod
	msgBody["content"] = r
	jsonData, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("failed to construct msgBody")
		return "failed to construct msgBody", err
	}
	err = Mq.Publish("minik8s", r.Spec.NodeName, "application/json", jsonData)
	if err != nil {
		fmt.Println("failed to send to mq")
		return "failed to send to mq", err
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

func AddServerlessFuncPod(funcName string) (string, error) {
	str, err := EtcdCli.Get(functionPrefix + funcName)
	if err != nil {
		fmt.Println(err)
		return err.Error(), err
	}

	var tmp function.Function
	err = json.Unmarshal(str, &tmp)
	if err != nil {
		fmt.Println(err)
		return err.Error(), err
	}

	imageName := tmp.Spec.ImageName

	id, _ := idgenerate.GenerateID()

	// 构建然后转json
	pod_ := &pod.Pod{}
	pod_.APIVersion = tmp.APIVersion
	pod_.Kind = "Pod"
	pod_.Status.Phase = pod.PodPending
	pod_.Metadata.Name = funcName + id[:8]
	pod_.Metadata.Namespace = DefaultNamespace
	pod_.Metadata.UUID = id
	pod_.Metadata.Labels = make(map[string]string)
	pod_.Metadata.Labels["serverlessFuncName"] = funcName
	pod_.Spec.Containers = make([]container.Container, 0)
	pod_.Spec.Containers = append(pod_.Spec.Containers, container.Container{
		Name:  pod_.Metadata.Name,
		Image: imageName,
	})

	value, err := json.Marshal(pod_)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 先查看一下key是否已经存在
	exist, err := EtcdCli.Exist(podPrefix + "Default" + "/" + pod_.Metadata.Name)
	if err != nil {
		fmt.Println("failed to check existence in etcd")
		return "failed to check existence in etcd", err
	}
	if exist {
		fmt.Println("pod has existed")
		return "pod has existed", nil
	}
	// 然后存入etcd
	err = EtcdCli.Put(podPrefix+"Default"+"/"+pod_.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	msgBody := make(map[string]interface{})
	msgBody["type"] = "create_pod"
	msgBody["content"] = pod_
	jsonData, err := json.Marshal(msgBody)
	if err != nil {
		fmt.Println("failed to construct msgBody")
		return "failed to construct msgBody", err
	}
	Mq.Publish("minik8s", "scheduler", "application/json", jsonData)

	return "add pod to minik8s", nil
}

func GetServerlessFuncPod(funcName string) ([]pod.Pod, error) {
	result := make([]pod.Pod, 0)
	pairs, err := EtcdCli.GetWithPrefix(podPrefix)
	if err != nil {
		return result, err
	}

	for _, p := range pairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		n, ok := tmp.Metadata.Labels["serverlessFuncName"]
		if ok && n == funcName {
			result = append(result, tmp)
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
