package worker

import (
	"encoding/json"
	"fmt"

	apiobject_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
)

var NodeRuntimeMangaer = minik8s_runtime.NewRuntimeManager()

type APIServer struct {
	IP   string
	Port string
}

type WorkerQueueBufferSize int

// Default buffer size for a pod worker's queue is 10.
const (
	WorkerQueueBufferSizeDefault WorkerQueueBufferSize = 10
)

// Every pod has a pod worker that is responsible for managing the pod.
type PodWorker struct {
	// API server IP and port.
	APIServer

	// Queue of tasks to be executed of one pod.
	TaskQueue chan *Task

	// AddTaskHandler func(*apiobject_pod.Pod) (string, error)
	// TODO: Implement the following handlers
	// UpdateTaskHandler func(*apiobject_pod.Pod) error
	// RemoveTaskHandler func(*apiobject_pod.Pod) error
}

func NewPodWorker(apiServer APIServer) *PodWorker {
	return &PodWorker{
		APIServer: apiServer,
		TaskQueue: make(chan *Task, WorkerQueueBufferSizeDefault),
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
		fmt.Println("PodWorker Add task")
		p.AddTaskHandler(task.Pod)
	case Task_Update:
	case Task_Remove:
	}
}

func (p *PodWorker) Stop() {
	close(p.TaskQueue)
}

func (p *PodWorker) AddTaskHandler(pod *apiobject_pod.Pod) (string, error) {
	_, err := NodeRuntimeMangaer.CreatePod(pod)
	if err != nil {
		return "", err
	}

	// Send http request to the API server to update the pod status.
	request_url := fmt.Sprintf("http://%s:%s/api/v1/UpdatePod", p.APIServer.IP, p.APIServer.Port)
	requestBody := make(map[string]interface{})
	requestBody["namespace"] = "default"
	requestBody["pod"] = pod
	requestBodyData, _ := json.Marshal(requestBody)
	response, err := httpRequest.PostRequest(
		request_url,
		requestBodyData,
	)
	if err != nil {
		fmt.Println("Error posting request: ", err)
		return "", err
	}
	fmt.Println("\nUpdate pod response: ", response)

	return pod.Metadata.UUID, nil
}

// TODO: Implement the following handlers
func (p *PodWorker) UpdateTaskHandler(pod *apiobject_pod.Pod) error {
	return nil
}

// TODO: Implement the following handlers
func (p *PodWorker) RemoveTaskHandler(pod *apiobject_pod.Pod) error {
	return nil
}
