package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/hpa"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

// 增删改查

func AddHPA(namespace string, desc *yaml.HPADesc) (string, error) {
	// 先检查是否存在
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(hpaPrefix + namespace + "/" + desc.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if existed {
		return "hpa existed", errors.New("hpa existed")
	}

	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	// 构建然后转json
	rs := &hpa.HPA{}
	rs.APIVersion = desc.ApiVersion
	rs.Kind = desc.Kind
	rs.Metadata.UUID = id
	rs.Metadata.Name = desc.Metadata.Name
	rs.Metadata.Namespace = namespace
	rs.Metadata.Labels = desc.Metadata.Labels
	rs.Spec = desc.Spec
	rs.Status = hpa.HPAStatus{}

	value, err := json.Marshal(rs)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	err = EtcdCli.Put(hpaPrefix+namespace+"/"+desc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add namespace to minik8s", nil
}

func RemoveHPA(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(hpaPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "hpa not found", nil
	}

	err = EtcdCli.Del(hpaPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func ModifyHPA(namespace string, name string) error {
	return nil
}

func GetHPA(namespace string, name string) (hpa.HPA, error) {
	mu.RLock()
	defer mu.RUnlock()
	var rs hpa.HPA

	v, err := EtcdCli.Get(hpaPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rs, err
	}

	err = json.Unmarshal(v, &rs)
	if err != nil {
		fmt.Println("failed to translate into json")
	}

	podPairs, err := EtcdCli.GetWithPrefix(podPrefix)
	if err != nil {
		fmt.Println("failed to get all pods from etcd")
		return rs, err
	}

	for _, pair := range podPairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(pair.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
		} else {
			value, ok := tmp.Metadata.Labels["hpa"]
			fmt.Println(tmp.Metadata.Labels, rs)
			if ok && value == rs.Metadata.Name {
				// TODO conditions
				rs.Status.ReadyReplicas += 1
			}
		}
	}

	return rs, nil
}

func GetHPAs(namespace string) (map[string]hpa.HPA, error) {
	mu.RLock()
	defer mu.RUnlock()
	rsmap := make(map[string]hpa.HPA, 0)

	pairs, err := EtcdCli.GetWithPrefix(hpaPrefix + namespace)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rsmap, err
	}

	for _, p := range pairs {
		var tmp hpa.HPA
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
		return rsmap, err
	}
	for _, pair := range podPairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(pair.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
		} else {
			value, ok := tmp.Metadata.Labels["hpa"]
			if ok {
				_, ok2 := rsmap[value]
				if ok2 {
					// TODO conditions
					tmp1 := rsmap[value]
					tmp1.Status.ReadyReplicas += 1
					rsmap[value] = tmp1
				}
			}
		}
	}

	return rsmap, nil
}

func GetAllHPAs() (map[string]hpa.HPA, error) {
	mu.RLock()
	defer mu.RUnlock()
	rsmap := make(map[string]hpa.HPA, 0)

	pairs, err := EtcdCli.GetWithPrefix(hpaPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rsmap, err
	}

	for _, p := range pairs {
		var tmp hpa.HPA
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
		return rsmap, err
	}
	for _, pair := range podPairs {
		var tmp pod.Pod
		err := json.Unmarshal([]byte(pair.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
		} else {
			value, ok := tmp.Metadata.Labels["hpa"]
			if ok {
				_, ok2 := rsmap[value]
				if ok2 {
					// TODO conditions
					tmp1 := rsmap[value]
					tmp1.Status.ReadyReplicas += 1
					rsmap[value] = tmp1
				}
			}
		}
	}

	return rsmap, nil
}
