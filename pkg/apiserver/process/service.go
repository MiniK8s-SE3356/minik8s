package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddService(namespace string, serviceType string, content string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(serviceType)
	if serviceType == "" {
		// 默认是ClusterIP
		serviceType = "ClusterIP"
	}

	if serviceType == "ClusterIP" {
		var desc yaml.ServiceClusterIPDesc
		err := json.Unmarshal([]byte(content), &desc)
		if err != nil {
			fmt.Println("failed to unmarshal")
			return "failed to unmarshal", err
		}
		clusterIP := service.ClusterIP{}
		id, _ := idgenerate.GenerateID()
		clusterIP.Metadata.Id = id
		clusterIP.ApiVersion = desc.ApiVersion
		clusterIP.Kind = desc.Kind
		clusterIP.Metadata.Namespace = namespace
		clusterIP.Metadata.Name = desc.Metadata.Name
		clusterIP.Metadata.Ip = ""
		clusterIP.Metadata.Labels = desc.Metadata.Labels
		clusterIP.Spec = desc.Spec
		clusterIP.Spec.Type = serviceType
		clusterIP.Status.Phase = service.CLUSTERIP_NOTREADY
		clusterIP.Status.Version = 0

		existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + clusterIP.Metadata.Name)
		if err != nil {
			return "failed to check existence in etcd", err
		}
		if existed {
			return "namespace existed", errors.New("namespace existed")
		}

		value, err := json.Marshal(clusterIP)
		if err != nil {
			fmt.Println("failed to marshal clusterIP")
			return "failed to marshal clusterIP", err
		}

		err = EtcdCli.Put(servicePrefix+namespace+"/"+clusterIP.Metadata.Name, string(value))
		if err != nil {
			fmt.Println("failed to write in etcd")
			return "failed to write in etcd", err
		}

		return "add successfully", nil

	} else if serviceType == "NodePort" {
		var desc yaml.ServiceNodePortDesc
		err := json.Unmarshal([]byte(content), &desc)
		if err != nil {
			fmt.Println("failed to unmarshal")
			return "failed to unmarshal", err
		}

		nodePort := service.NodePort{}
		id, _ := idgenerate.GenerateID()
		nodePort.Metadata.Id = id
		nodePort.Metadata.Name = desc.Metadata.Name
		nodePort.ApiVersion = desc.ApiVersion
		nodePort.Kind = desc.Kind
		nodePort.Metadata.Labels = desc.Metadata.Labels
		nodePort.Metadata.Namespace = namespace
		nodePort.Spec = desc.Spec
		nodePort.Spec.Type = serviceType
		nodePort.Status.Phase = service.NODEPORT_NOTREADY
		nodePort.Status.Version = 0

		existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + nodePort.Metadata.Name)
		if err != nil {
			return "failed to check existence in etcd", err
		}
		if existed {
			return "namespace existed", errors.New("namespace existed")
		}

		value, err := json.Marshal(nodePort)
		if err != nil {
			fmt.Println("failed to marshal nodePort")
			return "failed to marshal nodePort", err
		}
		fmt.Println(string(value))

		err = EtcdCli.Put(servicePrefix+namespace+"/"+nodePort.Metadata.Name, string(value))
		if err != nil {
			fmt.Println("failed to write in etcd")
			return "failed to write in etcd", err
		}

		return "add successfully", nil
	}

	return "", nil
}

func RemoveService(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "service not found", nil
	}

	content, _ := EtcdCli.Get(servicePrefix + namespace + "/" + name)
	var tmp map[string]interface{}
	err = json.Unmarshal(content, &tmp)
	if err != nil {
		fmt.Println("failed to unmarshal", err)
	}

	spec := tmp["spec"].(map[string]interface{})
	if spec["type"].(string) == "ClusterIP" {
		return removeClusterIP(namespace, name)
	} else if spec["type"].(string) == "NodePort" {
		return removeNodePort(namespace, name)
	} else {
		return "invalid service type", errors.New("invalid service type")
	}
}

func removeNodePort(namespace string, name string) (string, error) {
	var err error
	content, _ := EtcdCli.Get(servicePrefix + namespace + "/" + name)
	var np service.NodePort
	err = json.Unmarshal(content, &np)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return "failed to unmarshal", nil
	}
	cIPID := np.Status.ClusterIPID

	targetName := ""
	cIPsContent, _ := EtcdCli.GetWithPrefix(servicePrefix)
	for _, p := range cIPsContent {
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal", err)
		}

		spec := tmp["spec"].(map[string]interface{})
		if spec["type"].(string) == "ClusterIP" {
			var tmp service.ClusterIP
			err := json.Unmarshal([]byte(p.Value), &tmp)
			if err != nil {
				fmt.Println("failed to unmarshal")
				continue
			}

			if tmp.Metadata.Id == cIPID {
				targetName = tmp.Metadata.Name
				break
			}
		}
	}

	err = EtcdCli.Del(servicePrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}
	if targetName != "" {
		err = EtcdCli.Del(servicePrefix + namespace + "/" + targetName)
		if err != nil {
			fmt.Println("failed to del in etcd")
			return "failed to del in etcd", err
		}
	}

	return "del successfully", nil
}

func removeClusterIP(namespace string, name string) (string, error) {
	var err error
	content, _ := EtcdCli.Get(servicePrefix + namespace + "/" + name)
	var np service.NodePort
	err = json.Unmarshal(content, &np)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return "failed to unmarshal", nil
	}

	err = EtcdCli.Del(servicePrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func UpdateService(namespace string, name string, value string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	// existed, err := EtcdCli.Exist(servicePrefix + namespace + "/" + name)
	// if err != nil {
	// 	return "failed to check existence in etcd", err
	// }
	// if !existed {
	// 	return "service not found", errors.New("service not found")
	// }

	err := EtcdCli.Put(servicePrefix+namespace+"/"+name, value)
	if err != nil {
		fmt.Println("failed to update in etcd")
		return "failed to update in etcd", err
	}

	return "update successfully", nil
}

func GetService(namespace string, name string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	return "", nil
}

func GetServices(namespace string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	return "", nil
}

func GetAllService() (map[string][]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string][]interface{}, 2)
	pairs, err := EtcdCli.GetWithPrefix(servicePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	clusterIPArray := []interface{}{}
	nodePortArray := []interface{}{}
	for _, p := range pairs {
		// 先解析一下
		var tmp map[string]interface{}
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal", err)
		}

		spec := tmp["spec"].(map[string]interface{})
		if spec["type"].(string) == "ClusterIP" {
			var tmp service.ClusterIP
			err := json.Unmarshal([]byte(p.Value), &tmp)
			if err != nil {
				fmt.Println("failed to unmarshal")
			} else {
				clusterIPArray = append(clusterIPArray, tmp)
			}
		} else if spec["type"].(string) == "NodePort" {
			var tmp service.NodePort
			err := json.Unmarshal([]byte(p.Value), &tmp)
			if err != nil {
				fmt.Println("failed to unmarshal")
			} else {
				nodePortArray = append(nodePortArray, tmp)
			}
		} else {
			fmt.Println("invalid service type")
		}
	}

	result["clusterIP"] = clusterIPArray
	result["nodePort"] = nodePortArray

	// jsonData, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println("failed to translate into json")
	// 	return result, err
	// }
	return result, nil
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
