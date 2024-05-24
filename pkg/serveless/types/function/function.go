package function

import (
	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
)

type Function struct {
	apiobject.Basic `json:",inline" yaml:",inline"`
	Spec            FunctionSpec `json:"spec" yaml:"spec"`
}

type FunctionSpec struct {
	FileContent []byte `json:"fileContent" yaml:"fileContent"`
	FilePath    string `json:"filePath" yaml:"filePath"`
}
