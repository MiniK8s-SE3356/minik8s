package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddService(namespace string, desc *yaml.ServiceDesc) (string, error) {
	// serviceType, ok := desc.Spec["type"].(string)
	// fmt.Println(serviceType)
	// if !ok {
	// 	// 默认是ClusterIP
	// 	serviceType = "ClusterIP"
	// }

	// if serviceType == "ClusterIP" {
	// 	var desc yaml.ServiceClusterIPDesc
	// 	err := yaml.Unmarshal(b, &desc)
	// 	if err != nil {
	// 		fmt.Println("failed to unmarshal clusterIP service ", err.Error())
	// 		return err
	// 	}

	// 	requestMsg.Desc = desc
	// 	requestMsg.Namespace = namespace

	// 	jsonData, err = json.Marshal(requestMsg)
	// 	if err != nil {
	// 		fmt.Println("failed to marshal clusterIP service", err.Error())
	// 		return err
	// 	}
	// } else if serviceType == "NodePort" {
	// 	var desc yaml.ServiceNodePortDesc
	// 	err := yaml.Unmarshal(b, &desc)
	// 	if err != nil {
	// 		fmt.Println("failed to unmarshal NodePort service ", err.Error())
	// 		return err
	// 	}

	// 	requestMsg.Desc = desc
	// 	requestMsg.Namespace = namespace

	// 	jsonData, err = json.Marshal(requestMsg)
	// 	if err != nil {
	// 		fmt.Println("failed to marshal NodePort service", err.Error())
	// 		return err
	// 	}
	// }

	// existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + desc.Metadata.Name)
	// if err != nil {
	// 	return "failed to check existence in etcd", err
	// }
	// if existed {
	// 	return "namespace existed", errors.New("namespace existed")
	// }

	// err = EtcdCli.Put(servicePrefix+namespace+"/"+desc.Metadata.Name, value)
	// if err != nil {
	// 	fmt.Println("failed to write in etcd")
	// 	return "failed to write in etcd", err
	// }
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate id", err)
		return "", err
	}

	var a service.ClusterIP
	a.ApiVersion = "v1"
	a.Kind = "Service"
	a.Metadata.Name = "test-service"
	a.Metadata.Namespace = "Default"
	a.Metadata.Id = id
	a.Metadata.Ip = ""
	a.Spec.Ports = append(a.Spec.Ports, service.ClusterIPPortInfo{
		Protocal:   "TCP",
		Port:       1111,
		TargetPort: 1111,
	})
	a.Spec.Type = "ClusterIP"

	value, err := json.Marshal(a)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	fmt.Println(string(value))

	// existed, err := EtcdCli.Exist(servicePrefix + "Default" + "/" + a.Metadata.Name)
	// if err != nil {
	// 	return "failed to check existence in etcd", err
	// }
	// if existed {
	// 	return "namespace existed", errors.New("namespace existed")
	// }

	err = EtcdCli.Put(servicePrefix+namespace+"/"+a.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write in etcd")
		return "failed to write in etcd", err
	}

	return "add successfully", nil
}

func RemoveService(namespace string, name string) (string, error) {
	return "", nil
}

func UpdateService(namespace string, name string, value string) (string, error) {
	existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "namespace not found", errors.New("namespace not found")
	}

	err = EtcdCli.Put(servicePrefix+namespace+"/"+name, value)
	if err != nil {
		fmt.Println("failed to update in etcd")
		return "failed to update in etcd", err
	}

	return "update successfully", nil
}

func GetService(namespace string, name string) (string, error) {
	return "", nil
}

func GetServices(namespace string) (string, error) {
	return "", nil
}

func GetAllService() (string, error) {
	pairs, err := EtcdCli.GetWithPrefix(servicePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "", err
	}

	result := make(map[string][]etcdclient.KVPair, 2)
	clusterIPArray := []etcdclient.KVPair{}
	nodePortArray := []etcdclient.KVPair{}
	for _, p := range pairs {
		// 先解析一下
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal", err)
		}

		spec := tmp["spec"].(map[string]interface{})
		if spec["type"].(string) == "ClusterIP" {
			clusterIPArray = append(clusterIPArray, p)
		} else if spec["type"].(string) == "NodePort" {
			nodePortArray = append(nodePortArray, p)
		} else {
			fmt.Println("invalid service type")
		}
	}

	result["clusterIP"] = clusterIPArray
	result["nodePort"] = nodePortArray

	jsonData, err := json.Marshal(result)
	fmt.Println(jsonData)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	return string(jsonData), nil
}

func DescribeService(namespace string, name string) (string, error) {
	return "", nil
}

func DescribeServices(namespace string) (string, error) {
	return "", nil
}

func DescribeAllServices() (string, error) {
	return "", nil
}
