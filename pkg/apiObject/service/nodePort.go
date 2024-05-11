package service

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"

const (
	NODEPORT_PREFIX           = "MINIK8S-NODEPORT-"
	NODEPORT_ID_ALLOCATED     = "ID Allocated"
	NODEPORT_CLUSTERIP_FINISH = "Cluster Finish"
	NODEPORT_SUCCESS          = "Success"
	NODEPORT_ERROR            = "Error"
)

type NodePort struct {
	ApiVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   NodePortMetadata `json:"metadata" yaml:"metadata"`
	Spec       NodePortSpec     `json:"spec" yaml:"spec"`
	Status     NodePortStatus   `json:"status" yaml:"status"`
}

type NodePortMetadata struct {
	Name      string            `json:"name" yaml:"name"`
	Namespace string            `json:"namespace" yaml:"namespace"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	Id        string            `json:"id" yaml:"id"`
}

type NodePortSpec struct {
	Type     string             `json:"type" yaml:"type"`
	Selector selector.Selector  `json:"selector" yaml:"selector"`
	Ports    []NodePortPortInfo `json:"ports" yaml:"ports"`
}

type NodePortPortInfo struct {
	NodePort   int16  `json:"nodePort" yaml:"nodePort"`
	Protocal   string `json:"protocal" yaml:"protocal"`
	Port       int16  `json:"port" yaml:"port"`
	TargetPort int16  `json:"targetPort" yaml:"targetPort"`
}

type NodePortStatus struct {
	Phase       string `json:"phase" yaml:"phase"` /*只允许此文件上述的和状态有关的const常量，可参考飞书《Service设计方案》*/
	ClusterIPID string `json:"clusterIPID" yaml:"clusterIPID"`
}
