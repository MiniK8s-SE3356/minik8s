package worker

import apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"

type TaskType string

// TODO: add more task types
const (
	Task_Add    TaskType = "Add"
	Task_Remove TaskType = "Remove"
	Task_Update TaskType = "Update"
)

type Task struct {
	Type TaskType
	Pod  *apiobject.Pod
}
