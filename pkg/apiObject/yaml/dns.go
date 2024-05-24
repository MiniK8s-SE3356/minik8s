package yaml

import (
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
)

type DNSDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`
	Spec dns.DnsSpec `yaml:"spec" json:"spec"`
}

// type DNSSpec struct {
// 	Host  string `yaml:"host" json:"host"`
// 	Paths []struct {
// 		SubPath string `yaml:"subPath" json:"subPath"`
// 		SvcName string `yaml:"svcName" json:"svcName"`
// 		SvcPort uint16 `yaml:"svcPort" json:"svcPort"`
// 	} `yaml:"paths" json:"paths"`
// }
