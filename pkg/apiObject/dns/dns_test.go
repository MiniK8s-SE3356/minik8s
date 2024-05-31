package dns_test

import (
	"os"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
	"github.com/google/uuid"
)

var dns_object dns.Dns

func TestMain(m *testing.M) {
	// pre-test code
	dns_object.ApiVersion = "v1"
	dns_object.Kind = "Dns"
	dns_object.Metadata = dns.DnsMetadata{
		Name:      "test-dns1",
		Namespace: "default",
		Id:        uuid.NewString(),
	}
	dns_object.Spec = dns.DnsSpec{
		Host:  "dns.test",
		Paths: []dns.DnsPathInfo{},
	}
	dns_object.Status = dns.DnsStatus{
		Phase:       dns.DNS_NOTREADY,
		Version:     0,
		PathsStatus: make(map[string]dns.DnsPathStatus),
	}

	// test func
	exitCode := m.Run()

	// post-test code

	// 返回测试运行的退出码
	os.Exit(exitCode)
}

func TestDns(t *testing.T) {
	if dns_object.ApiVersion != "v1" || dns_object.Kind != "Dns" || dns_object.Metadata.Name != "test-dns1" || dns_object.Spec.Host != "dns.test" || dns_object.Status.Phase != dns.DNS_NOTREADY {
		t.Fatalf("Dns test error, can't match the init field value.\n")
	}
}
