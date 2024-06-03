package persistVolume

const (
	PV_PHASE_CREATED   = "Created"
	PV_PHASE_AVAILABLE = "Available"
	PV_PHASE_RELEASED  = "Released"
	PV_PHASE_FAILED    = "Failed"
	PV_NAME_PREFIX     = "MINIK8S-PV-"
	PV_TYPE_NFS        = "nfs"
)

type PersistVolume struct {
	ApiVersion string                `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                `json:"kind" yaml:"kind"` /*只允许PersistVolume*/
	Metadata   PersistVolumeMetadata `json:"metadata" yaml:"metadata"`
	Spec       PersistVolumeSpec     `json:"spec" yaml:"spec"`
	Status     PersistVolumeStatus   `json:"status" yaml:"status"`
}

type PersistVolumeMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Id     string            `json:"id" yaml:"id"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type PersistVolumeSpec struct {
	Type     string `json:"type" yaml:"type"` /*只允许nfs*/
	Capacity string `json:"capacity" yaml:"capacity"`
}

type PersistVolumeStatus struct {
	Phase    string   `json:"phase" yaml:"phase"`
	MountPod []string `json:"mountPod" yaml:"mountPod"`
}
