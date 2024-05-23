package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

func AddNode(desc *node.Node) (node.Node, error) {
	mu.Lock()
	defer mu.Unlock()
	// 检查name是否为空
	result := node.Node{}
	if desc.Metadata.Name != "" {
		return result, nil
	}

	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate id")
		return result, err
	}
	desc.Metadata.Id = id
	desc.Metadata.Name = desc.Status.Hostname + "@" + desc.Status.Ip
	// desc.Metadata.Labels = map[string]string{}
	// desc.Status.Condition = []string{}

	// 检查是否已经存在
	existed, err := EtcdCli.Exist(nodePrefix + desc.Metadata.Name)
	if err != nil {
		return result, err
	}

	if existed {
		return result, nil
	}

	value, err := json.Marshal(desc)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return result, err
	}

	// 然后存入etcd
	// TODO node的name还需要处理
	err = EtcdCli.Put(nodePrefix+desc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return result, err
	}

	result = *desc
	return result, nil
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

func GetNode(name string) (node.Node, error) {
	mu.RLock()
	defer mu.RUnlock()
	var result node.Node
	existed, err := EtcdCli.Exist(nodePrefix + name)
	if err != nil {
		return result, err
	}
	if !existed {
		return result, nil
	}

	v, err := EtcdCli.Get(nodePrefix + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}
	err = json.Unmarshal(v, &result)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return result, err
	}

	return result, nil
}

func GetNodes() ([]node.Node, error) {
	mu.RLock()
	defer mu.RUnlock()
	var result []node.Node
	pairs, err := EtcdCli.GetWithPrefix(nodePrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	for _, p := range pairs {
		var n node.Node
		err := json.Unmarshal([]byte(p.Value), &n)
		if err != nil {
			fmt.Println("failed to unmarshal")
		} else {
			result = append(result, n)
		}
	}

	return result, nil
}

func DescribeNode(name string) (string, error) {
	return "", nil
}

func DescribeAllNodes() (string, error) {
	return "", nil
}

func NodeHeartBeat(nodeStatus node.NodeStatus, pods []pod.Pod, nodePorts []service.NodePort) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	name := nodeStatus.Hostname + "@" + nodeStatus.Ip

	nodeInfo, err := EtcdCli.Get(nodePrefix + name)
	if err != nil || len(nodeInfo) == 0 {
		fmt.Println("failed to get node from etcd")
		return "failed to get node from etcd", err
	}
	var tmp node.Node
	err = json.Unmarshal(nodeInfo, &tmp)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return "failed to unmarshal", err
	}
	tmp.Status = nodeStatus
	nodeInfo, err = json.Marshal(tmp)
	if err != nil {
		fmt.Println("failed to marshal")
		return "failed to marshal", err
	}
	err = EtcdCli.Put(nodePrefix+name, string(nodeInfo))
	if err != nil {
		fmt.Println("failed to write node to etcd")
		return "failed to write node to etcd", err
	}

	for _, p := range pods {
		ns := DefaultNamespace
		name := p.Metadata.Name
		podInfo, err := EtcdCli.Get(podPrefix + ns + "/" + name)
		if err != nil || len(podInfo) == 0 {
			fmt.Println("failed to get node from etcd")
			continue
			// return "failed to get node from etcd", err
		}
		var tmp pod.Pod
		err = json.Unmarshal(podInfo, &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
			continue
			// return "failed to unmarshal", err
		}
		tmp.Status = p.Status
		podInfo, err = json.Marshal(tmp)
		if err != nil {
			fmt.Println("failed to marshal")
			continue
			// return "failed to marshal", err
		}
		err = EtcdCli.Put(podPrefix+ns+"/"+name, string(podInfo))
		if err != nil {
			fmt.Println("failed to write node to etcd")
			// return "failed to write pod to etcd", err
		}

	}

	return "", nil
}
