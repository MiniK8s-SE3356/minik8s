package hpa

import (
	"time"

	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
)

type HPA struct {
	apiobject.Basic `yaml:",inline"`
	Spec            yaml.HPASpec `yaml:"spec" json:"spec"`
	Status          HPAStatus    `yaml:"status" json:"status"`
}

type HPAStatus struct {
	ReadyReplicas  int       `json:"readyReplicas" yaml:"readyReplicas"`
	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	CPUUsage       float64   `yaml:"CPUUsage" json:"CPUUsage"`
	MemUsage       float64   `yaml:"MemUsage" json:"MemUsage"`
}

// type HPACondition struct {
// 	Type           string    `json:"type" yaml:"type"`
// 	Status         string    `json:"status" yaml:"status"`
// 	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
// 	Reason         string    `json:"reason" yaml:"reason"`
// 	Message        string    `json:"message" yaml:"message"`
// }
