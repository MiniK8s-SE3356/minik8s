package service_test

import (
	"os"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	// pre-test code

	// test func
	exitCode := m.Run()

	// post-test code

	os.Exit(exitCode)
}

const (
	test_label = "[NodePort2ClusterIP Test Error]"
)

func TestNodePort2ClusterIP(t *testing.T) {
	nodeport := service.NodePort{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: service.NodePortMetadata{
			Name:      "test-nodeport1",
			Namespace: "default",
			Labels:    map[string]string{},
			Id:        uuid.NewString(),
		},
		Spec: service.NodePortSpec{
			Type:     "NodePort",
			Selector: selector.Selector{MatchLabels: map[string]string{}},
			Ports:    []service.NodePortPortInfo{},
		},
		Status: service.NodePortStatus{
			Phase:       service.NODEPORT_NOTREADY,
			Version:     0,
			ClusterIPID: "",
		},
	}
	nodeport.Metadata.Labels["test1"] = "test1"
	nodeport.Spec.Selector.MatchLabels["test2"] = "test2"
	nodeport.Spec.Ports = append(nodeport.Spec.Ports, service.NodePortPortInfo{
		NodePort:   uint16(41234),
		Protocol:   "tcp",
		Port:       uint16(80),
		TargetPort: uint16(8080),
	})

	clusterip := service.NodePort2ClusterIP(&nodeport)

	if clusterip.Kind != "Service" || clusterip.Metadata.Labels["test1"] != "test1" || clusterip.Spec.Type != "ClusterIP" || clusterip.Spec.Selector.MatchLabels["test2"] != "test2" {
		t.Fatalf("%s metadata error\n", test_label)
	}
	nlength := len(nodeport.Spec.Ports)
	clenght := len(clusterip.Spec.Ports)
	if nlength != clenght {
		t.Fatalf("%s port lenght does not match\n", test_label)
	}

	for i := 1; i < nlength; i++ {
		if nodeport.Spec.Ports[i].TargetPort != clusterip.Spec.Ports[i].TargetPort || nodeport.Spec.Ports[i].Port != clusterip.Spec.Ports[i].Port || nodeport.Spec.Ports[i].Protocol != clusterip.Spec.Ports[i].Protocol {
			t.Fatalf("%s port infomation does not match\n", test_label)
		}
	}
}
