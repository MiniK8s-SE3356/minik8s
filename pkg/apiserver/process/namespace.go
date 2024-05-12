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
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	id = idgenerate.NamespacePrefix + id

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
	err = EtcdCli.Put(namespacePrefix+id, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add namespace to minik8s", nil
}

func RemoveNamespace(name string) (string, error) {
	pairs, err := EtcdCli.GetWithPrefix(namespacePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "", err
	}

	var target string
	for _, p := range pairs {
		var tmp namespace.Namespace
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal namespace")
			return "", err
		}

		if tmp.Metadata.Name == name {
			target = p.Key
			break
		}
	}

	if target == "" {
		fmt.Println("namespace not found")
		return "namespace not found", nil
	}

	err = EtcdCli.Del(target)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

// func GetNamespace(name string) (string, error) {
// 	return "", nil
// }

func GetNamespaces() (string, error) {
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
