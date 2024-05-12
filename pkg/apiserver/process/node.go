package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddNode(desc *node.Node) (string, error) {
	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	id = node.NODE_PREFIX + id

	desc.Metadata.Id = id

	value, err := json.Marshal(desc)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 然后存入etcd
	err = EtcdCli.Put(nodePrefix+id, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add node to minik8s", nil
}

func RemoveNode(name string) (string, error) {
	pairs, err := EtcdCli.GetWithPrefix(nodePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "", err
	}

	var target string
	for _, p := range pairs {
		var tmp node.Node
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal node")
			return "", err
		}

		if tmp.Metadata.Name == name {
			target = p.Key
			break
		}
	}

	if target == "" {
		fmt.Println("node not found")
		return "node not found", nil
	}

	err = EtcdCli.Del(target)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil

}

func ModifyNode() (string, error) {
	return "", nil
}

func GetNode(name string) (string, error) {
	pairs, err := EtcdCli.GetWithPrefix(nodePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return "", err
	}

	var target string
	for _, p := range pairs {
		var tmp node.Node
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal node")
			return "", err
		}

		if tmp.Metadata.Name == name {
			target = p.Value
			break
		}
	}

	if target == "" {
		fmt.Println("node not found")
		return "node not found", nil
	}

	return target, nil
}

func GetNodes() (string, error) {
	pairs, err := EtcdCli.GetWithPrefix(nodePrefix)
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

func DescribeNode(name string) (string, error) {
	return "", nil
}

func DescribeAllNodes() (string, error) {
	return "", nil
}
