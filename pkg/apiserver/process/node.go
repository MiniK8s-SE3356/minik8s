package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
)

func AddNode(desc *node.Node) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	// 检查name是否为空
	if desc.Metadata.Name == "" {
		return "empty name is not allowed", nil
	}

	// 检查是否已经存在
	existed, err := EtcdCli.Exist(nodePrefix + desc.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}

	if existed {
		return "node has existed", nil
	}

	value, err := json.Marshal(desc)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	// 然后存入etcd
	// TODO node的name还需要处理
	err = EtcdCli.Put(nodePrefix+desc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add node to minik8s", nil
}

func RemoveNode(name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(nodePrefix + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "node not found", nil
	}

	err = EtcdCli.Del(nodePrefix + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func ModifyNode() (string, error) {
	mu.Lock()
	defer mu.Unlock()
	return "", nil
}

func GetNode(name string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()
	existed, err := EtcdCli.Exist(nodePrefix + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "node not found", nil
	}

	// TODO

	return "", nil
}

func GetNodes() (string, error) {
	mu.RLock()
	defer mu.RUnlock()
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

func NodeHeartBeat() {

}
