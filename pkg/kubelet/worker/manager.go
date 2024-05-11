package worker

import (
	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
)

// PodManager will receive pod requests from the kubelet and manage pod workers to handle them.
type PodManager interface {
	AddPod(pod *apiobject.Pod) error
	UpdatePod(pod *apiobject.Pod) error
	RemovePod(pod *apiobject.Pod) error
	GetPods() ([]*apiobject.Pod, error)
	GetPodByUID(uid string) (*apiobject.Pod, error)
	GetPodByName(namespace, name string) (*apiobject.Pod, error)
}

// podManager is the default implementation of PodManager.
type podManager struct {
	// podWorkers is a map of pod workers indexed by pod UID.
	PodWorkers map[string]PodWorker
}

func NewPodManager() PodManager {
	return &podManager{
		PodWorkers: make(map[string]PodWorker),
	}
}

func (pm *podManager) AddPod(pod *apiobject.Pod) error {
	// TODO: Implement this method
	return nil
}

func (pm *podManager) UpdatePod(pod *apiobject.Pod) error {
	return nil
}

func (pm *podManager) RemovePod(pod *apiobject.Pod) error {
	return nil
}

func (pm *podManager) GetPods() ([]*apiobject.Pod, error) {
	return nil, nil
}

func (pm *podManager) GetPodByUID(uid string) (*apiobject.Pod, error) {
	return nil, nil
}

func (pm *podManager) GetPodByName(namespace, name string) (*apiobject.Pod, error) {
	return nil, nil
}
