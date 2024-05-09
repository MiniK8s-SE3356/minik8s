package runtime

import (
	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	runtime_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types"
)

type RuntimeManager struct {
	containerManager *runtime_container.ContainerManager
	imageManager     *image.ImageManager
}

func (rm *RuntimeManager) CreatePod(pod *apiobject.Pod) (string, error) {
	pauseContainerId, err := rm.CreateAndStartPauseContainer(pod)
	if err != nil {
		return "", err
	}

	// Create containers for the pod
	for _, container := range pod.Spec.Containers {
		containerConfig := &types.CreateContainerConfig{
			Image:       container.Image,
			NetworkMode: "container:" + pauseContainerId,
			IpcMode:     "container:" + pauseContainerId,
			PidMode:     "container:" + pauseContainerId,
		}
		containerId, err := rm.containerManager.CreateContainer(container.Name, containerConfig)
		if err != nil {
			return "", err
		}
		// Start the container
		_, err = rm.containerManager.StartContainer(containerId)
		if err != nil {
			return "", err
		}
	}

	return pod.Metadata.UUID, nil
}

func (rm *RuntimeManager) CreateAndStartPauseContainer(pod *apiobject.Pod) (string, error) {
	// First, try to pull pause container's image
	_, err := rm.imageManager.PullImage("registry.aliyuncs.com/google_containers/pause:3.9")
	if err != nil {
		return "", err
	}

	uuid := pod.Metadata.UUID
	pauseName := "pause-" + uuid

	// Create a container for pause
	// TODO: Set the pause container's network
	pauseContainerId, err := rm.containerManager.CreateContainer(pauseName, &types.CreateContainerConfig{
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
