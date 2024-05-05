package ty

type ServiceDesc struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name   string            `json:"name" yaml:"name"`
		Labels map[string]string `json:"labels" yaml:"labels"`
	} `json:"metadata" yaml:"metadata"`
}

type ServiceSpec struct {
	Selector  map[string]string `json:"selector" yaml:"selector"`
	Ports     []ServicePortDesc `json:"ports" yaml:"ports"`
	Type      string            `json:"type" yaml:"type"`
	ClusterIP string            `json:"clusterIP" yaml:"clusterIP"`
}

type ServicePortDesc struct {
	Name       string `json:"name" yaml:"name"`
	Port       int    `json:"port" yaml:"port"`
	NodePort   int    `json:"nodePort" yaml:"nodePort"`
	TargetPort int    `json:"targetPort" yaml:"targetPort"`
	Protocol   string `json:"protocol" yaml:"protocol"`
}
