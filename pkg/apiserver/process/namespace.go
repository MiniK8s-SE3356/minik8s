package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/namespace"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

// TODO

func AddNamespace(desc *yaml.NamespaceDesc) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	// 检查name是否为空
	if desc.Metadata.Name == "" {
		return "empty name is not allowed", nil
	}

	// 检查是否已经存在
	existed, err := EtcdCli.Exist(namespacePrefix + desc.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}

	if existed {
		return "namespace has existed", nil
	}

	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	// 构建然后转json
	namespace := &namespace.Namespace{}
	namespace.ApiVersion = desc.ApiVersion
	namespace.Kind = desc.Kind
	namespace.Metadata.Name = desc.Metadata.Name
	namespace.Metadata.Id = id
	namespace.Metadata.Labels = desc.Metadata.Labels

	value, err := json.Marshal(namespace)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 然后存入etcd
	err = EtcdCli.Put(namespacePrefix+namespace.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add namespace to minik8s", nil
}

func RemoveNamespace(name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	// TODO: 删除namespace下的其他资源
	existed, err := EtcdCli.Exist(namespacePrefix + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "namespace not found", nil
	}

	err = EtcdCli.Del(namespacePrefix + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func GetNamespace(name string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	result, err := EtcdCli.Get(namespacePrefix + name)
	if err != nil {
		fmt.Println("failed to get namespace from etcd")
		return "", err
	}

	return string(result), nil
}

func GetNamespaces() (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	pairs, err := EtcdCli.GetWithPrefix(namespacePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "", err
	}

	jsonData, err := json.Marshal(pairs)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	return string(jsonData), nil
}
