package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
)

func GetAllEndpoint() (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	pairs, err := EtcdCli.GetWithPrefix(endpointPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return map[string]interface{}{}, nil
	}

	result := make(map[string]interface{})

	for _, p := range pairs {
		var tmp service.EndPoint
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal endpoint")
		} else {
			result[p.Key] = tmp
		}
	}

	return result, nil
}

func UpdateEndpointBatch(desc map[string]service.EndPoint) (map[string]interface{}, error) {
	mu.Lock()
	defer mu.Unlock()
	result := make(map[string]interface{})

	pairs, err := EtcdCli.GetWithPrefix(endpointPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}
	for _, p := range pairs {
		err := EtcdCli.Del(p.Key)
		if err != nil {
			fmt.Println("failed to del in etcd")
		}
	}

	jsonDataArr := make(map[string]string)
	for k, v := range desc {
		jsonData, err := json.Marshal(v)
		if err != nil {
			fmt.Println("failed to write in etcd")
			return result, err
		}

		jsonDataArr[k] = string(jsonData)
	}

	for k, v := range jsonDataArr {
		err := EtcdCli.Put(k, v)
		if err != nil {
			fmt.Println("failed to write in etcd")
		}
	}

	return result, nil
}
