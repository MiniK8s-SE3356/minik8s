package persistVolume

import "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"

const (
	PVC_PHASE_AVAILABLE = "Available"
	PVC_PHASE_FAILED    = "Failed"
	PVC_NAME_PREFIX     = "MINIK8S-PVC-"
)

type PersistVolumeClaim struct {
	ApiVersion string                     `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                     `json:"kind" yaml:"kind"` /*只允许PersistVolumeClaim*/
	Metadata   PersistVolumeClaimMetadata `json:"metadata" yaml:"metadata"`
	Spec       PersistVolumeClaimSpec     `json:"spec" yaml:"spec"`
	Status     PersistVolumeClaimStatus   `json:"status" yaml:"status"`
}

type PersistVolumeClaimMetadata struct {
	Name string `json:"name" yaml:"name"`
	Id string `json:"id" yaml:"id"`
}

type PersistVolumeClaimSpec struct {
	Type     string            `json:"type" yaml:"type"` /*只允许nfs*/
	Selector selector.Selector `json:"selector" yaml:"selector"`
}

type PersistVolumeClaimStatus struct {
	Phase   string   `json:"phase" yaml:"phase"`
	BoundPV []string `json:"boundPV" yaml:"boundPV"`
}
