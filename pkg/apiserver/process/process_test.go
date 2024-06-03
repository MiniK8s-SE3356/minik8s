package process_test

import (
	"encoding/json"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
)

func TestPod(t *testing.T) {
	pairs, _ := process.EtcdCli.GetWithPrefix("/mink8s/pod/Default/test")
	if len(pairs) != 0 {
		panic("not empty")
	}

	var p yaml.PodDesc
	p.ApiVersion = "v1"
	p.Kind = "Pod"
	p.Metadata.Name = "test"
	p.Spec.Containers = make([]container.Container, 0)
	var c container.Container
	c.Image = "docker.io/library/nginx"
	var po container.ContainerPort
	po.ContainerPort = 80
	po.HostPort = 80
	c.Ports = append(c.Ports, po)
	p.Spec.Containers = append(p.Spec.Containers, c)

	_, err := process.AddPod("Default", &p)
	if err != nil {
		panic(err.Error())
	}

	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/pod/Default/")
	if len(pairs) != 1 {
		panic("not added")
	}

	process.RemovePod("Default", "test")
	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/pod/Default/")
	if len(pairs) != 0 {
		panic("not removed")
	}
}

func TestNode(t *testing.T) {
	pairs, _ := process.EtcdCli.GetWithPrefix("/mink8s/node/host@1.1.1.1")
	if len(pairs) != 0 {
		panic("not empty")
	}

	var n node.Node
	n.Status.Hostname = "host"
	n.Status.Ip = "1.1.1.1"
	_, err := process.AddNode(&n)
	if err != nil {
		panic(err.Error())
	}

	pairs, _ = process.EtcdCli.GetWithPrefix("/mink8s/node")
	if len(pairs) != 1 {
		panic("not added")
	}

	process.RemoveNode("host@1.1.1.1")
	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/node")
	if len(pairs) != 0 {
		panic("not removed")
	}
}

func TestService(T *testing.T) {
	pairs, _ := process.EtcdCli.GetWithPrefix("/mink8s/service/Default/test")
	if len(pairs) != 0 {
		panic("not empty")
	}

	var s yaml.ServiceDesc
	s.ApiVersion = "v1"
	s.Kind = "Service"
	s.Metadata.Name = "test"
	jsonData, _ := json.Marshal(s)

	_, err := process.AddService("Default", "ClusterIP", string(jsonData))
	if err != nil {
		panic(err.Error())
	}

	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/service/Default/")
	if len(pairs) != 1 {
		panic("not added")
	}

	process.RemoveService("Default", "test")
	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/service/Default/")
	if len(pairs) != 0 {
		panic("not removed")
	}
}

func TestReplicaset(T *testing.T) {
	pairs, _ := process.EtcdCli.GetWithPrefix("/mink8s/replicaset/Default/test")
	if len(pairs) != 0 {
		panic("not empty")
	}

	var s yaml.ReplicaSetDesc
	s.ApiVersion = "v1"
	s.Kind = "Replicaset"
	s.Metadata.Name = "test"

	_, err := process.AddReplicaSet("Default", &s)
	if err != nil {
		panic(err.Error())
	}

	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/replicaset/Default/")
	if len(pairs) != 1 {
		panic("not added")
	}

	process.RemoveReplicaSet("Default", "test")
	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/replicaset/Default/")
	if len(pairs) != 0 {
		panic("not removed")
	}
}

func TestHPA(T *testing.T) {
	pairs, _ := process.EtcdCli.GetWithPrefix("/mink8s/hpa/Default/test")
	if len(pairs) != 0 {
		panic("not empty")
	}

	var s yaml.HPADesc
	s.ApiVersion = "v1"
	s.Kind = "HPA"
	s.Metadata.Name = "test"

	_, err := process.AddHPA("Default", &s)
	if err != nil {
		panic(err.Error())
	}

	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/hpa/Default/")
	if len(pairs) != 1 {
		panic("not added")
	}

	process.RemoveHPA("Default", "test")
	pairs, _ = process.EtcdCli.GetWithPrefix("/minik8s/hpa/Default/")
	if len(pairs) != 0 {
		panic("not removed")
	}
}
