package worker

import (
	"encoding/json"
	"fmt"
	"sync"

	apiobject_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	minik8s_runtime "github.com/MiniK8s-SE3356/minik8s/pkg/runtime"
	cadvisorutils "github.com/MiniK8s-SE3356/minik8s/pkg/utils/cAdvisorUtils"
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

	// Local pod object.
	LocalPod apiobject_pod.Pod

	// worker lock
	Lock sync.RWMutex
}

func NewPodWorker(apiServer APIServer) *PodWorker {
	return &PodWorker{
		APIServer: apiServer,
		TaskQueue: make(chan *Task, WorkerQueueBufferSizeDefault),
	}
}

func (p *PodWorker) AddTask(task *Task) error {
	// TODO: check queue size is not full
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
		fmt.Println("PodWorker Remove task")
		p.RemoveTaskHandler(task.Pod, task.Callback)
	case Task_Restart:
		fmt.Println("PodWorker Restart task")
		p.RestartTaskHandler(task.Pod)
	}
}

func (p *PodWorker) Stop() {
	close(p.TaskQueue)
}

func (p *PodWorker) AddTaskHandler(pod *apiobject_pod.Pod) (string, error) {
	// Lock the worker
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, err := NodeRuntimeMangaer.CreatePod(pod)
	if err != nil {
		fmt.Println("Add Task Error!")
		return "", err
	}

	// Set local pod
	p.LocalPod = *pod

	// Send http request to the API server to update the pod status.
	request_url := fmt.Sprintf("http://%s:%s/api/v1/UpdatePod", p.APIServer.IP, p.APIServer.Port)
	requestBody := make(map[string]interface{})
	requestBody["namespace"] = "default"
	requestBody["pod"] = pod
	requestBodyData, _ := json.Marshal(requestBody)

	fmt.Println("Update pod request: ", string(requestBodyData))

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
	p.Lock.Lock()
	defer p.Lock.Unlock()
	return nil
}

// TODO: Implement the following handlers
func (p *PodWorker) RemoveTaskHandler(pod *apiobject_pod.Pod, callback func(arg interface{})) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	// If remove pod success, send http request to the API server to update the pod status.
	// and stop the pod worker. Remove the pod worker from the pod manager map.

	err := NodeRuntimeMangaer.RemovePod(pod)
	if err != nil {
		fmt.Println("Remove Task Error!")
		return err
	}

	// call the callback function to remove podworker from the pod manager map
	callback(pod.Metadata.UUID)
	return nil
}

func (p *PodWorker) RestartTaskHandler(pod *apiobject_pod.Pod) error {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, err := NodeRuntimeMangaer.RestartPod(pod)
	if err != nil {
		fmt.Println("Restart Task Error!")
		return err
	}

	p.LocalPod = *pod

	request_url := fmt.Sprintf("http://%s:%s/api/v1/UpdatePod", p.APIServer.IP, p.APIServer.Port)
	requestBody := make(map[string]interface{})
	requestBody["namespace"] = "default"
	requestBody["pod"] = pod
	requestBodyData, _ := json.Marshal(requestBody)

	fmt.Println("Update pod request: ", string(requestBodyData))

	response, err := httpRequest.PostRequest(
		request_url,
		requestBodyData,
	)
	if err != nil {
		fmt.Println("Error posting request: ", err)
		return err
	}
	fmt.Println("\nUpdate pod response: ", response)

	return nil
}

func (p *PodWorker) FetchandUpdateLocalPod() {
	// Try to lock pod worker's read lock, if the pod worker is locked by other writer
	// we give up this updating operation. Otherwise we access pod's status.
	// We don't want pod's writer's operation blocking the heartbeat sending!
	isLocked := p.Lock.TryRLock()
	if !isLocked {
		return
	}
	defer p.Lock.RUnlock()

	// calculate cpu and memory usage for one pod
	podCPUUsage := 0.0
	podMemUsage := 0.0

	for _, container := range p.LocalPod.Spec.Containers {
		cpuUsage, memoryUsage, err := cadvisorutils.GetContainerCPUandMemory(
			"localhost",
			// TODO: change port
			"8081",
			container.Name,
		)
		if err != nil {
			// This means a container in the pod is not running.
			p.LocalPod.Status.Phase = apiobject_pod.PodFailed
			// Add a task to pod worker to restart this pod.
			p.AddTask(&Task{
				Type: Task_Restart,
				Pod:  &p.LocalPod,
			})
		}
		podCPUUsage += cpuUsage
		podMemUsage += memoryUsage
	}

	p.LocalPod.Status.CPUUsage = podCPUUsage
	p.LocalPod.Status.MemoryUsage = podMemUsage
}
