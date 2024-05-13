package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddPod(desc *yaml.PodDesc) (string, error) {
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
	pod.Metadata.UUID = id
	pod.Metadata.Labels = desc.Metadata.Labels
	// for _, c := range pod.Spec.Containers {
	// 	// var tmp container.Container
	// 	pod.Spec.Containers = append(pod.Spec.Containers, c)
	// }

	value, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 先查看一下key是否已经存在
	EtcdCli.Get(pod.Metadata.Name)
	// 然后存入etcd
	err = EtcdCli.Put(podPrefix+id, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add pod to minik8s", nil
}

func RemovePod(pod string, name string) (string, error) {
	return "", nil
}

func ModifyPod() (string, error) {
	return "", nil
}

func GetPod(pod string, name string) (string, error) {
	return "", nil
}

func GetPods(pod string) (string, error) {
	return "", nil
}

func GetAllPods() (string, error) {
	return "", nil
}

func DescribePod(pod string, name string) (string, error) {
	return "", nil
}

func DescribePods(pod string) (string, error) {
	return "", nil
}

func DescribeAllPods() (string, error) {
	return "", nil
}
