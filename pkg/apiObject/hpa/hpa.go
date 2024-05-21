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
	Replicas   int            `json:"replicas" yaml:"replicas"`
	Conditions []HPACondition `json:"conditions" yaml:"conditions"`
}

type HPACondition struct {
	Type           string    `json:"type" yaml:"type"`
	Status         string    `json:"status" yaml:"status"`
	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	Reason         string    `json:"reason" yaml:"reason"`
	Message        string    `json:"message" yaml:"message"`
}
