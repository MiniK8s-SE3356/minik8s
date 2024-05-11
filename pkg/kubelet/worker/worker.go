package worker

import (
	apiobject_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
)

var NodeRuntimeMangaer = minik8s_runtime.NewRuntimeManager()

type WorkerQueueBufferSize int

// Default buffer size for a pod worker's queue is 10.
const (
	WorkerQueueBufferSizeDefault WorkerQueueBufferSize = 10
)

// Every pod has a pod worker that is responsible for managing the pod.
type PodWorker struct {
	// Queue of tasks to be executed of one pod.
	TaskQueue chan *Task

	AddTaskHandler func(*apiobject_pod.Pod) (string, error)
	// TODO: Implement the following handlers
	UpdateTaskHandler func(*apiobject_pod.Pod) error
	RemoveTaskHandler func(*apiobject_pod.Pod) error
}

func NewPodWorker() *PodWorker {
	return &PodWorker{
		TaskQueue: make(chan *Task, WorkerQueueBufferSizeDefault),

		AddTaskHandler: NodeRuntimeMangaer.CreatePod,
	}
}

func (p *PodWorker) AddTask(task *Task) error {
	p.TaskQueue <- task
	return nil
}

func (p *PodWorker) Run() {
	for task := range p.TaskQueue {
		p.ExecuteTask(task)
	}
}

func (p *PodWorker) ExecuteTask(task *Task) {
	switch task.Type {
	case Task_Add:
		p.AddTaskHandler(task.Pod)
	case Task_Update:
	case Task_Remove:
	}
}

func (p *PodWorker) Stop() {
	close(p.TaskQueue)
}
