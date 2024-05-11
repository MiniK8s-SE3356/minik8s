package namespace

const (
	NAMESPACE_PERFIX = "MINIK8S-NAMESPACE-"
)

type Namespace struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Kind       string            `json:"kind" yaml:"kind"`
	Metadata   NamespaceMetadata `json:"metadata" yaml:"metadata"`
}

type NamespaceMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
	Id     string            `json:"id" yaml:"id"`
}
