package dns

const (
	DNS_PREFIX   = "MINIK8S-PREFIX-"
	DNS_READY    = "READY"
	DNS_NOTREADY = "NOTREADY"
)

type Dns struct {
	ApiVersion string      `json:"apiVersion" yaml:"apiVersion"`
	Kind       string      `json:"kind" yaml:"kind"`
	Metadata   DnsMetadata `json:"metadata" yaml:"metadata"`
	Spec       DnsSpec     `json:"spec" yaml:"spec"`
	Status     DnsStatus   `json:"status" yaml:"status"`
}

type DnsMetadata struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Id        string `json:"id" yaml:"id"`
}

type DnsSpec struct {
	Host  string        `json:"host" yaml:"host"`
	Paths []DnsPathInfo `json:"paths" yaml:"paths"`
}

type DnsStatus struct {
	Phase   string `json:"phase" yaml:"phase"` /*READY or NOTREADY*/
	Version int    `json:"version" yaml:"version"`
	/* PathsStatus的key为对应的svcName*/
	PathsStatus map[string]DnsPathStatus `json:"pathsStatus" yaml:"pathsStatus"`
}

type DnsPathInfo struct {
	SubPath string `json:"subPath" yaml:"subPath"`
	SvcName string `json:"svcName" yaml:"svcName"`
	SvcPort uint16 `json:"svcPort" yaml:"svcPort"`
}

// func (di *DnsPathInfo)SpecEuqalStatus(ds *DnsPathStatus){
// 	return (di.SubPath==ds.SubPath)&&(di.)
// }

type DnsPathStatus struct {
	SubPath string `json:"subPath" yaml:"subPath"`
	SvcIP   string `json:"svcIP" yaml:"svcIP"`
	SvcPort uint16 `json:"svcPort" yaml:"svcPort"`
}
