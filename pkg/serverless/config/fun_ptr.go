package config

import "github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/workflow"

var GetFuncionPodRequestFrequency func()map[string]float64
var TriggerServerlessFunction func(string, string, string) (string, error) 
var TriggerServerlessWorkflow func(workflow.Workflow,string) 
