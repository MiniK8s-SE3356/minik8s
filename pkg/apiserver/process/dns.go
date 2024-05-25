package process

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
)

// 增删改查

func AddDNS(namespace string, desc *yaml.DNSDesc) (string, error) {
	// 先检查是否存在
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(dnsPrefix + namespace + "/" + desc.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if existed {
		return "dns existed", errors.New("dns existed")
	}

	id, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate uuid")
		return "failed to generate uuid", err
	}

	// 构建然后转json
	rs := &dns.Dns{}
	rs.ApiVersion = desc.ApiVersion
	rs.Kind = desc.Kind
	rs.Metadata.Id = id
	rs.Metadata.Name = desc.Metadata.Name
	rs.Metadata.Namespace = namespace
	// rs.Metadata.Labels = desc.Metadata.Labels
	rs.Spec = desc.Spec
	rs.Status = dns.DnsStatus{}
	rs.Status.Phase = dns.DNS_NOTREADY

	value, err := json.Marshal(rs)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return "failed to translate into json ", err
	}

	err = EtcdCli.Put(dnsPrefix+namespace+"/"+desc.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "add namespace to minik8s", nil
}

func RemoveDNS(namespace string, name string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(dnsPrefix + namespace + "/" + name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "dns not found", nil
	}

	err = EtcdCli.Del(dnsPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to del in etcd")
		return "failed to del in etcd", err
	}

	return "del successfully", nil
}

func UpdateDNS(namespace string, dns dns.Dns) (string, error) {
	// 先检查是否存在
	mu.Lock()
	defer mu.Unlock()
	existed, err := EtcdCli.Exist(dnsPrefix + namespace + "/" + dns.Metadata.Name)
	if err != nil {
		return "failed to check existence in etcd", err
	}
	if !existed {
		return "dns not exist", errors.New("dns existed")
	}

	value, err := json.Marshal(dns)
	if err != nil {
		fmt.Println("failed to marshal dns")
		return "failed to marshal dns", err
	}

	err = EtcdCli.Put(dnsPrefix+namespace+"/"+dns.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return "failed to write to etcd", err
	}

	return "update dns success", nil
}

func UpdateDNSList(dnsList map[string]dns.Dns) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	for n, d := range dnsList {
		existed, err := EtcdCli.Exist(dnsPrefix + DefaultNamespace + "/" + n)
		if err != nil {
			fmt.Println("failed to check existence in etcd", err)
			continue
		}
		if !existed {
			fmt.Println("dns not exist")
			continue
		}

		value, err := json.Marshal(d)
		if err != nil {
			fmt.Println("failed to marshal dns")
			continue
		}

		err = EtcdCli.Put(dnsPrefix+DefaultNamespace+"/"+d.Metadata.Name, string(value))
		if err != nil {
			fmt.Println("failed to write to etcd ", err.Error())
		}
	}
	return "update dns success", nil
}

func GetDNS(namespace string, name string) (dns.Dns, error) {
	mu.RLock()
	defer mu.RUnlock()
	var rs dns.Dns

	v, err := EtcdCli.Get(dnsPrefix + namespace + "/" + name)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rs, err
	}

	err = json.Unmarshal(v, &rs)
	if err != nil {
		fmt.Println("failed to translate into json")
	}

	return rs, nil
}

func GetDNSs(namespace string) (map[string]dns.Dns, error) {
	mu.RLock()
	defer mu.RUnlock()
	rsmap := make(map[string]dns.Dns, 0)

	pairs, err := EtcdCli.GetWithPrefix(dnsPrefix + namespace)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rsmap, err
	}

	for _, p := range pairs {
		var tmp dns.Dns
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			rsmap[tmp.Metadata.Name] = tmp
		}
	}

	return rsmap, nil
}

func GetAllDNSs() ([]dns.Dns, error) {
	mu.RLock()
	defer mu.RUnlock()
	rsmap := make([]dns.Dns, 0)

	pairs, err := EtcdCli.GetWithPrefix(dnsPrefix)
	if err != nil {
		fmt.Println("failed to get from etcd")
		return rsmap, err
	}

	for _, p := range pairs {
		var tmp dns.Dns
		err := json.Unmarshal([]byte(p.Value), &tmp)
		if err != nil {
			fmt.Println("failed to translate into json")
		} else {
			rsmap = append(rsmap, tmp)
		}
	}

	return rsmap, nil
}
