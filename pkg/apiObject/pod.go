package apiobject

import "github.com/MiniK8s-SE3356/minik8s/pkg/types"

type HostPath struct {
	Path string `json:"path" yaml:"path"`

	Type string `json:"type" yaml:"type"`
}

type Volume struct {
	Name string `json:"name" yaml:"name"`

	HostPath HostPath `json:"hostPath" yaml:"hostPath"`
}

type PodSpec struct {
	NodeName string `json:"nodeName" yaml:"nodeName"`

	Containers []types.Container `json:"containers" yaml:"containers"`

	Volumes []Volume `json:"volumes" yaml:"volumes"`

	// Some labels for the pod
	NodeSelector map[string]string `json:"nodeSelector" yaml:"nodeSelector"`
}

type Pod struct {
	Basic `yaml:",inline"`

	Spec PodSpec `json:"spec" yaml:"spec"`

	Status types.PodStatus `json:"status" yaml:"status"`
}
