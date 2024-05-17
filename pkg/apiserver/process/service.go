package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
)

func AddService(namespace string, desc interface{}) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	serviceType, ok := desc.(map[string]interface{})["Spec"].(map[string]interface{})["type"].(string)
	fmt.Println(serviceType)
	if !ok {
		// 默认是ClusterIP
		serviceType = "ClusterIP"
	}

	if serviceType == "ClusterIP" {
		clusterIP := desc.(service.ClusterIP)
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
		nodePort := desc.(service.NodePort)
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
				clusterIPArray = append(clusterIPArray, tmp)
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
