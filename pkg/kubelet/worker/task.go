package worker

import apiobject_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"

type TaskType string

// TODO: add more task types
const (
	Task_Add    TaskType = "Add"
	Task_Remove TaskType = "Remove"
	Task_Update TaskType = "Update"
)

type Task struct {
	Type TaskType
	Pod  *apiobject_pod.Pod
}
