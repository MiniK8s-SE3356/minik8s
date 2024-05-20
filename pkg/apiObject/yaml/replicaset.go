package yaml

import "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"

type ReplicaSetDesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`
	Spec ReplicaSetSpec `yaml:"spec" json:"spec"`
}

type ReplicaSetSpec struct {
	Replicas int `yaml:"replicas" json:"replicas"`
	Selector struct {
		MatchLabels map[string]string `yaml:"matchLabels" json:"matchLabels"`
	} `yaml:"selector" json:"selector"`
	Template struct {
		Metadata struct {
			Name   string            `yaml:"name" json:"name"`
			Labels map[string]string `yaml:"labels" json:"labels"`
		} `yaml:"metadata" json:"metadata"`
		Spec struct {
			Containers []container.Container `yaml:"containers" json:"containers"`
		} `yaml:"spec" json:"spec"`
	}
}
