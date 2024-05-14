package replicaset

import (
	"time"

	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
)

type Replicaset struct {
	apiobject.Basic `yaml:",inline"`
	Spec            yaml.ReplicaSetSpec `yaml:"spec" json:"spec"`
	Status          ReplicaSetStatus    `yaml:"status" json:"status"`
}

type ReplicaSetStatus struct {
	Replicas      int                   `json:"replicas" yaml:"replicas"`
	ReadyReplicas int                   `json:"readyReplicas" yaml:"readyReplicas"`
	Conditions    []ReplicaSetCondition `json:"conditions" yaml:"conditions"`
}

type ReplicaSetCondition struct {
	Type           string    `json:"type" yaml:"type"`
	Status         string    `json:"status" yaml:"status"`
	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	Reason         string    `json:"reason" yaml:"reason"`
	Message        string    `json:"message" yaml:"message"`
}
