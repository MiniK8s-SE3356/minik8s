package dns

import (
	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
)

const (
	DNSReady    string = "Ready"
	DNSNotReady string = "NotReady"
)

type DNS struct {
	apiobject.Basic `yaml:",inline"`
	Spec            yaml.DNSSpec `yaml:"spec" json:"spec"`
	Status          DNSStatus    `yaml:"status" json:"status"`
}

type DNSStatus struct {
	Phase       string `yaml:"phase" json:"phase"`
	Version     int    `yaml:"version" json:"version"`
	PathsStatus map[string]struct {
		SubPath string `yaml:"subPath" json:"subPath"`
		SvcIP   string `yaml:"svcIP" json:"svcIP"`
		SvcPort int    `yaml:"svcPort" json:"svcPort"`
	} `yaml:"pathsStatus" json:"pathsStatus"`
}
