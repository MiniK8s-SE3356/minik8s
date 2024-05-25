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
			result[tmp.Id] = tmp
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

func AddorDeleteEndpoint(d []string, a []service.EndPoint) {
	mu.Lock()
	defer mu.Unlock()
	for _, name := range d {
		existed, err := EtcdCli.Exist(endpointPrefix + name)
		if err != nil {
			fmt.Println("failed to check existence in etcd", err)
		}
		if !existed {
			fmt.Println("endpoint not found")
		}

		err = EtcdCli.Del(endpointPrefix + name)
		if err != nil {
			fmt.Println("failed to del in etcd")
		}
	}

	for _, ep := range a {
		// 检查name是否为空
		if ep.Id == "" {
			fmt.Println("empty name is not allowed")
			continue
		}

		value, err := json.Marshal(ep)
		if err != nil {
			fmt.Println("failed to translate into json ", err.Error())
		}

		// 然后存入etcd
		err = EtcdCli.Put(endpointPrefix+ep.Id, string(value))
		if err != nil {
			fmt.Println("failed to write to etcd ", err.Error())
		}
	}

}
