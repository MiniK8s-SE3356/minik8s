package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/replicaset"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

// 增删改查

func AddReplicaSet(namespace string, desc *yaml.ReplicaSetDesc) (string, error) {
	// 先检查是否存在
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(replicasetPrefix + namespace + "/" + desc.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if existed {
		return "replicaset existed", errors.New("replicaset existed")
	}

	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	// 构建然后转json
	rs := &replicaset.Replicaset{}
	rs.APIVersion = desc.ApiVersion
	rs.Kind = desc.Kind
	rs.Metadata.UUID = id
	rs.Metadata.Name = desc.Metadata.Name
	rs.Metadata.Namespace = namespace
	rs.Metadata.Labels = desc.Metadata.Labels
	rs.Spec = desc.Spec
	rs.Status = replicaset.ReplicaSetStatus{}
	rs.Status.Conditions = []replicaset.ReplicaSetCondition{}

	value, err := json.Marshal(rs)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	err = EtcdCli.Put(replicasetPrefix+namespace+"/"+desc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add namespace to minik8s", nil
}

func RemoveReplicaSet(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(replicasetPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "replicaset not found", nil
	}

	err = EtcdCli.Del(replicasetPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func ModifyReplicaSet(namespace string, name string) error {
	return nil
}

func GetReplicaSet(namespace string, name string) (replicaset.Replicaset, error) {
	mu.RLock()
	defer mu.RUnlock()
	var result replicaset.Replicaset

	existed, err := EtcdCli.Exist(replicasetPrefix + namespace + "/" + name)
	if err != nil {
		return result, err
	}
	if !existed {
		return result, nil
	}

	tmp, err := EtcdCli.Get(replicasetPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, nil
	}

	err = json.Unmarshal(tmp, &result)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return result, nil
	}

	return result, nil
}

func GetReplicaSets(namespace string) ([]replicaset.Replicaset, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]replicaset.Replicaset, 0)

	pairs, err := EtcdCli.GetWithPrefix(replicasetPrefix + namespace)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp replicaset.Replicaset
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			result = append(result, tmp)
		}
	}

	return result, nil
}

func GetAllReplicaSets() ([]replicaset.Replicaset, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]replicaset.Replicaset, 0)

	pairs, err := EtcdCli.GetWithPrefix(replicasetPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return result, err
	}

	for _, p := range pairs {
		var tmp replicaset.Replicaset
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			result = append(result, tmp)
		}
	}

	return result, nil
}
