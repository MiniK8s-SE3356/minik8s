package process

import (
	"encoding/json"
	"errors"
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

	// TODO: pod转的时候需要更改，需要查看ectd中的持久化卷状态
	// 如果POD需要PVC而PVC不存在，直接error
	// 如果POD需要PV而PV不存在，直接error
	// 如果POD需要PVC,PVC存在，而PVC对应的PV不存在，立刻创建PV
	// 如果POD需要PVC,且PVC下辖所需PV，则正常绑定
	for _, volume_item := range desc.Spec.Volumes {
		// 此处，PV的优先级高于PVC,如果两个都写，以PV的规则为准
		if volume_item.HostPath.Path != "" {
			continue
		}
		// 这俩的优先级均低于HostPath
		if volume_item.PersistentVolume.PvName != "" {
			result_str := IsPVAvailable(DefaultNamespace, volume_item.PersistentVolume.PvName)
			if result_str != IsPVAvailable_Return_OK {
				// PV不是OK直接爆
				return result_str, errors.New(result_str)
			}
			continue
		}
		if volume_item.PersistentVolumeClaim.ClaimName != "" {
			result_str, spectype := IsPVC_PV_exist(DefaultNamespace, volume_item.PersistentVolumeClaim.ClaimName, volume_item.Name)
			if result_str == IsPVC_PV_exist_Return_OK {
				continue
			}
			if result_str == IsPVC_PV_exist_Return_PVMISS {
				err := AddPVImmediately(volume_item.Name, spectype)
				if err != nil {
					// 创建失败，也报错
					return err.Error(), err
				}
				continue
			}
			// 除了上述俩，直接爆
			return result_str, errors.New(result_str)
		}
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
