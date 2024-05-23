package yaml

import (
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
)

type HPADesc struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   struct {
		Name   string            `yaml:"name" json:"name"`
		Labels map[string]string `yaml:"labels" json:"labels"`
	} `yaml:"metadata" json:"metadata"`
	Spec HPASpec `yaml:"spec" json:"spec"`
}

type HPASpec struct {
	MinReplicas int `yaml:"minReplicas" json:"minReplicas"`
	MaxReplicas int `yaml:"maxReplicas" json:"maxReplicas"`
	Selector    struct {
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
	TimeInterval time.Duration `yaml:"timeInterval" json:"timeInterval"`
	Metrics      HPAMetrics    `yaml:"metrics" json:"metrics"`
}

type HPAMetrics struct {
	CPUPercent float64 `yaml:"cpuPercent" json:"cpuPercent"`
	MemPercent float64 `yaml:"memPercent" json:"memPercent"`
}
