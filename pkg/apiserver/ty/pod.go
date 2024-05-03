package ty

type PodDesc struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   struct {
		Name   string            `json:"name" yaml:"name"`
		Labels map[string]string `json:"labels" yaml:"labels"`
	} `json:"metadata" yaml:"metadata"`
	Spec []ContainerDesc `json:"spec" yaml:"spec"`
}

type ContainerDesc struct {
	Name      string     `json:"name" yaml:"name"`
	Image     string     `json:"image" yaml:"image"`
	Ports     []PortDesc `json:"ports" yaml:"ports"`
	Resources struct {
		Requests ResourceDesc `json:"requests" yaml:"requests"`
		Limits   ResourceDesc `json:"limits" yaml:"limits"`
	} `json:"resources" yaml:"resources"`
}

type PortDesc struct {
	ContainerPort int    `json:"containerPort" yaml:"containerPort"`
	HostPort      int    `json:"hostPort" yaml:"hostPort"`
	Protocol      string `json:"protocol" yaml:"protocol"`
	Name          string `json:"name" yaml:"name"`
	HostIP        string `json:"hostIP" yaml:"hostIP"`
}

type ResourceDesc struct {
	Memory string `json:"memory" yaml:"memory"`
	CPU    string `json:"cpu" yaml:"cpu"`
}
