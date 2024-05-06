package ty

type ServiceDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`

	// 这里因为spec的具体格式和service的type有关，所以这里是map[string]string的类型
	// 使用的时候先访问Spec["type"]获取类型，然后使用相应的类型再次进行解析
	Spec map[string]string `yaml:"spec" json:"spec"`
}

type ServiceClusterIPDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`

	Spec ServiceClusterIPSpec `yaml:"spec" json:"spec"`
}

type ServiceNodePortDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`

	Spec ServiceNodePortSpec `yaml:"spec" json:"spec"`
}

type ServiceClusterIPSpec struct {
	Type string `yaml:"type" json:"type"`

	Selector map[string]string `yaml:"selector" json:"selector"`
	Ports    []struct {
		Name       string `yaml:"name" json:"name"`
		Port       int    `yaml:"port" json:"port"`
		TargetPort int    `yaml:"targetPort" json:"targetPort"`
		Protocol   string `yaml:"protocol" json:"protocol"`
	} `yaml:"ports" json:"ports"`
	ClusterIP string `yaml:"clusterIP" json:"clusterIP"`
}

type ServiceNodePortSpec struct {
	Type string `yaml:"type" json:"type"`

	Selector map[string]string `yaml:"selector" json:"selector"`
	Ports    []struct {
		Name       string `yaml:"name" json:"name"`
		Port       int    `yaml:"port" json:"port"`
		NodePort   int    `yaml:"nodePort" json:"nodePort"`
		TargetPort int    `yaml:"targetPort" json:"targetPort"`
		Protocol   string `yaml:"protocol" json:"protocol"`
	} `yaml:"ports" json:"ports"`
	ClusterIP string `yaml:"clusterIP" json:"clusterIP"`
}
