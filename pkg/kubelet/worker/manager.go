package worker

import (
	"fmt"

	apiobject_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
)

// PodManager will receive pod requests from the kubelet and manage pod workers to handle them.
type PodManager interface {
	AddPod(pod *apiobject_pod.Pod) error
	UpdatePod(pod *apiobject_pod.Pod) error
	RemovePod(pod *apiobject_pod.Pod) error
	GetPods() ([]*apiobject_pod.Pod, error)
	GetPodByUID(uid string) (*apiobject_pod.Pod, error)
	GetPodByName(namespace, name string) (*apiobject_pod.Pod, error)
	GetPodNum() int
	FetchLocalPods() ([]apiobject_pod.Pod, error)
}

// podManager is the default implementation of PodManager.
type podManager struct {
	APIServer
	// podWorkers is a map of pod workers indexed by pod UID.
	PodWorkers map[string]*PodWorker
}

func NewPodManager(apiServer APIServer) PodManager {
	return &podManager{
		APIServer:  apiServer,
		PodWorkers: make(map[string]*PodWorker),
	}
}

func (pm *podManager) AddPod(pod *apiobject_pod.Pod) error {
	fmt.Println(pod.Metadata.Name + " is added to the pod manager.")
	UID := pod.Metadata.UUID
	fmt.Println("UID: " + UID)

	// Check if the pod worker already exists
	if _, ok := pm.PodWorkers[UID]; ok {
		fmt.Println("pod worker with UID " + UID + " already exists.")
		return fmt.Errorf("pod worker with UID %s already exists", UID)
	}

	// Create a new pod worker
	podWorker := NewPodWorker(pm.APIServer)
	pm.PodWorkers[UID] = podWorker

	// Create a go routine to run the pod worker
	go podWorker.Run()

	// Create add task for the pod worker
	addTask := &Task{
		Type: Task_Add,
		Pod:  pod,
	}

	// Add the task to the pod worker's queue
	err := podWorker.AddTask(addTask)
	if err != nil {
		return err
	}

	return nil
}

func (pm *podManager) UpdatePod(pod *apiobject_pod.Pod) error {
	fmt.Println("Updating pod: ", pod.Metadata.Name, " with UID: ", pod.Metadata.UUID)
	return nil
}

func (pm *podManager) RemovePod(pod *apiobject_pod.Pod) error {
	fmt.Println("Removing pod: ", pod.Metadata.Name, " with UID: ", pod.Metadata.UUID)
	UID := pod.Metadata.UUID

	if _, ok := pm.PodWorkers[UID]; !ok {
		fmt.Println("pod worker with UID " + UID + " does not exist.")
		return fmt.Errorf("pod worker with UID %s does not exist", UID)
	}

	removeTask := &Task{
		Type: Task_Remove,
		Pod:  pod,
	}

	err := pm.PodWorkers[UID].AddTask(removeTask)
	if err != nil {
		return err
	}

	return nil
}

func (pm *podManager) GetPods() ([]*apiobject_pod.Pod, error) {
	return nil, nil
}

func (pm *podManager) GetPodByUID(uid string) (*apiobject_pod.Pod, error) {
	return nil, nil
}

func (pm *podManager) GetPodByName(namespace, name string) (*apiobject_pod.Pod, error) {
	return nil, nil
}

func (pm *podManager) GetPodNum() int {
	return len(pm.PodWorkers)
}

func (pm *podManager) FetchLocalPods() ([]apiobject_pod.Pod, error) {
	result := make([]apiobject_pod.Pod, 0)
	for _, podWorker := range pm.PodWorkers {
		fmt.Println("Fetching local pods")
		podWorker.FetchandUpdateLocalPod()
		result = append(result, podWorker.LocalPod)
	}
	return result, nil
}
