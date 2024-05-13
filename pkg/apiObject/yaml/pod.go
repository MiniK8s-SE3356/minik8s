package yaml

// for pod yaml
type PodDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`
	Spec []ContainerDesc `yaml:"spec" json:"spec"`
}

type ContainerDesc struct {
	Name      string     `yaml:"name" json:"name"`
	Image     string     `yaml:"image" json:"image"`
	Ports     []PortDesc `yaml:"ports" json:"ports"`
	Resources struct {
		Requests ResourceDesc `yaml:"requests" json:"requests"`
		Limits   ResourceDesc `yaml:"limits" json:"limits"`
	} `yaml:"resources" json:"resources"`
}

type PortDesc struct {
	ContainerPort int    `yaml:"containerPort" json:"containerPort"`
	HostPort      int    `yaml:"hostPort" json:"hostPort"`
	Protocol      string `yaml:"protocol" json:"protocol"`
	Name          string `yaml:"name" json:"name"`
	HostIP        string `yaml:"hostIP" json:"hostIP"`
}

type ResourceDesc struct {
	Memory string `yaml:"memory" json:"memory"`
	CPU    string `yaml:"cpu" json:"cpu"`
}
