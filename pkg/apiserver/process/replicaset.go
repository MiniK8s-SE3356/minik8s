package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
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

func GetReplicaSet(namespace string, name string) (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	var r replicaset.Replicaset
	result := make(map[string]interface{}, 0)

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

	err = json.Unmarshal(tmp, &r)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return result, nil
	}

	result[replicasetPrefix+namespace+"/"+name] = r

	return result, nil
}

func GetReplicaSets(namespace string) ([]replicaset.Replicaset, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]replicaset.Replicaset, 0)
	rsmap := make(map[string]replicaset.Replicaset, 0)

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
			rsmap[tmp.Metadata.Name] = tmp
		}
	}

	podPairs, err := EtcdCli.GetWithPrefix(podPrefix)
	if err != nil {
		fmt.Println("failed to get all pods from etcd")
		return result, err
	}
	for _, pair := range podPairs {
		tmp := pod.Pod{}
		err := json.Unmarshal([]byte(pair.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
		} else {
			value, ok := tmp.Metadata.Labels["replicaset"]
			if ok {
				_, ok2 := rsmap[value]
				if ok2 {
					// TODO conditions
					tmp1 := rsmap[value]
					tmp1.Status.Replicas += 1
					rsmap[value] = tmp1
				}
			}
		}
	}
	for _, rs := range rsmap {
		result = append(result, rs)
	}

	return result, nil
}

func GetAllReplicaSets() (map[string]interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	result := make(map[string]interface{}, 0)

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
			result[p.Key] = tmp
		}
	}

	return result, nil
}
