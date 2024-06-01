package types

import apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"

type SlurmJob struct {
	apiobject.Basic `json:",inline" yaml:",inline"`
	JobID           string `json:"jobID" yaml:"jobID"`
	Partition       string `json:"partition" yaml:"partition"`
	Nodes           string `json:"nodes" yaml:"nodes"`
	NTasksPerNode   string `json:"ntasksPerNode" yaml:"ntasksPerNode"`
	CPUPerTask      string `json:"cpuPerTask" yaml:"cpuPerTask"`
	GPUNum          string `json:"gpuNum" yaml:"gpuNum"`
	OutputFile      string `json:"outputFile" yaml:"outputFile"`
	ErrFile         string `json:"errFile" yaml:"errFile"`

	// Runtime related
	// default to ~/minik8s/UUID
	WorkDir     string   `json:"workDir" yaml:"workDir"`
	CompileCmds []string `json:"compileCmds" yaml:"compileCmds"`
	Modules     []string `json:"modules" yaml:"modules"`
	RunCmds     []string `json:"runCmds" yaml:"runCmds"`

	// Username and password to login to the HPC platform
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	State string `json:"state" yaml:"state"`
}
