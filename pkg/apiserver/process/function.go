package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
)

func GetFunction(name string) ([]function.Function, error) {
	var result []function.Function
	var r function.Function

	existed, err := EtcdCli.Exist(functionPrefix + name)
	if err != nil {
		return result, err
	}
	if !existed {
		return result, nil
	}

	tmp, err := EtcdCli.Get(podPrefix + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, nil
	}

	err = json.Unmarshal(tmp, &r)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return result, nil
	}

	result = append(result, r)

	return result, nil
}

func GetAllFunction() ([]function.Function, error) {
	var funcs []function.Function

	pairs, err := EtcdCli.GetWithPrefix(functionPrefix)
	if err != nil {
		fmt.Println(err)
		return funcs, err
	}

	for _, v := range pairs {
		var tmp function.Function
		err := json.Unmarshal([]byte(v.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
			continue
		}

		funcs = append(funcs, tmp)
	}

	return funcs, nil
}

func GetAllServerlessFunction() ([]string, error) {
	var names []string

	pairs, err := EtcdCli.GetWithPrefix(functionPrefix)
	if err != nil {
		fmt.Println(err)
		return names, err
	}

	for _, v := range pairs {
		var tmp function.Function
		err := json.Unmarshal([]byte(v.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
			continue
		}

		names = append(names, tmp.Metadata.Name)
	}

	return names, nil
}

func RemoveFunction(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(functionPrefix + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "function not found", nil
	}

	err = EtcdCli.Del(functionPrefix + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil

}
