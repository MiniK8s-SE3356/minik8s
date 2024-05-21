package runtime

import (
	"fmt"

	minik8s_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	runtime_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
)

type RuntimeManager struct {
	containerManager *runtime_container.ContainerManager
	imageManager     *image.ImageManager
}

func NewRuntimeManager() *RuntimeManager {
	return &RuntimeManager{
		containerManager: &runtime_container.ContainerManager{},
		imageManager:     &image.ImageManager{},
	}
}

func (rm *RuntimeManager) CreatePod(pod *minik8s_pod.Pod) (string, error) {
	//!debug//
	// fmt.Println("Creating pod for:", pod.Metadata.UUID)
	//!debug//
	pauseContainerId, err := rm.CreateAndStartPauseContainer(pod)
	if err != nil {
		return "", err
	}
	//!debug//
	// fmt.Println("Pause container created")
	// fmt.Println("Container number: ", len(pod.Spec.Containers))
	//!debug//

	// Create containers for the pod
	for _, container := range pod.Spec.Containers {
		//!debug//
		// fmt.Println("Creating container: ", container.Image)
		//!debug//
		containerConfig := &minik8s_container.CreateContainerConfig{
			Image:       container.Image,
			NetworkMode: "container:" + pauseContainerId,
			IpcMode:     "container:" + pauseContainerId,
			PidMode:     "container:" + pauseContainerId,
		}
		//!debug//
		// fmt.Println("Start creating container")
		//!debug//
		containerId, err := rm.containerManager.CreateContainer(container.Name, containerConfig)
		//!debug//
		// fmt.Println("Finish creating container")
		//!debug//
		if err != nil {
			fmt.Println("Failed to create container", err)
			return "", err
		}
		// Start the container
		//!debug//
		// fmt.Println("Start starting container")
		//!debug//
		_, err = rm.containerManager.StartContainer(containerId)
		//!debug//
		// fmt.Println("Finish starting container")
		//!debug//
		if err != nil {
			fmt.Println("Failed to start container", err)
			return "", err
		}
	}

	return pod.Metadata.UUID, nil
}

func (rm *RuntimeManager) CreateAndStartPauseContainer(pod *minik8s_pod.Pod) (string, error) {
	// First, try to pull pause container's image
	_, err := rm.imageManager.PullImage("registry.aliyuncs.com/google_containers/pause:3.9")
	if err != nil {
		return "", err
	}

	uuid := pod.Metadata.UUID
	pauseName := "pause-" + uuid

	// Create a container for pause
	// TODO: Set the pause container's network
	pauseContainerId, err := rm.containerManager.CreateContainer(pauseName, &minik8s_container.CreateContainerConfig{
		Image:   "registry.aliyuncs.com/google_containers/pause:3.9",
		IpcMode: "shareable",
	})

	if err != nil {
		println("Failed to create pause container")
		return "", err
	}

	_, err = rm.containerManager.StartContainer(pauseContainerId)

	if err != nil {
		println("Failed to start pause container")
		return "", err
	}

	return pauseContainerId, nil
}
