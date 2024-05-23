package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	minik8s_pod "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	runtime_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	runtime_image "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	minik8s_node "github.com/MiniK8s-SE3356/minik8s/pkg/types/node"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/nettools"
	nodestatusutils "github.com/MiniK8s-SE3356/minik8s/pkg/utils/nodeStatusUtils"
	"github.com/docker/go-connections/nat"
)

type RuntimeManager struct {
	containerManager *runtime_container.ContainerManager
	imageManager     *runtime_image.ImageManager
}

func NewRuntimeManager() *RuntimeManager {
	return &RuntimeManager{
		containerManager: &runtime_container.ContainerManager{},
		imageManager:     &runtime_image.ImageManager{},
	}
}

func (rm *RuntimeManager) GetVolumeBinds(volumes *[]minik8s_pod.Volume, volumeMounts *[]minik8s_container.VolumeMount) []string {
	volumes_map := make(map[string]*minik8s_pod.Volume)

	// Make a map of volumes
	for _, volume := range *volumes {
		volumes_map[volume.Name] = &volume
	}

	// Create a list of volume binds
	// TODO: Check if the hostPath valid
	volumeBinds := []string{}
	for _, volumeMount := range *volumeMounts {
		if _, ok := volumes_map[volumeMount.Name]; !ok {
			fmt.Println("Volume ", volumeMount.Name, " not found")
			continue
		}

		volume := volumes_map[volumeMount.Name]

		volumeBind := fmt.Sprintf("%s:%s", volume.HostPath.Path, volumeMount.MountPath)

		volumeBinds = append(volumeBinds, volumeBind)
	}

	return volumeBinds
}

func (rm *RuntimeManager) CreatePod(pod *minik8s_pod.Pod) (string, error) {
	fmt.Println("Creating pod for ", pod.Metadata.Name, " with UUID ", pod.Metadata.UUID)
	jsonPod, _ := json.Marshal(pod)
	fmt.Println(string(jsonPod))

	// Create a pause container for the pod and this function will set pod's IP address
	pauseContainerId, err := rm.CreateAndStartPauseContainer(pod)
	if err != nil {
		return "", err
	}

	// Create containers for the pod
	for i, container := range pod.Spec.Containers {
		// parse environment variables
		env := []string{}
		for _, envVar := range container.Env {
			env = append(env, fmt.Sprintf("%s=%s", envVar.Name, envVar.Value))
		}
		// TODO: parse command

		// parse volume mounts. We should mount the pod's volumes to the container
		volumeBinds := rm.GetVolumeBinds(&pod.Spec.Volumes, &container.VolumeMounts)
		for _, volumeBind := range volumeBinds {
			fmt.Println("Volume bind: ", volumeBind)
		}

		// Because all containers in a pod share the same network namespace,
		// we have set the port bindings in the pause container. No need to set every container's port bindings.

		containerConfig := &minik8s_container.CreateContainerConfig{
			Image:       container.Image,
			NetworkMode: "container:" + pauseContainerId,
			IpcMode:     "container:" + pauseContainerId,
			PidMode:     "container:" + pauseContainerId,
			Env:         env,
			Binds:       volumeBinds,
			CPU:         container.Resources.Limits.CPU,
			Memory:      container.Resources.Limits.Memory,
		}
		containerId, err := rm.containerManager.CreateContainer(container.Name, containerConfig)
		pod.Spec.Containers[i].Id = containerId
		if err != nil {
			fmt.Println("Failed to create container", err)
			return "", err
		}
		// Start the container
		_, err = rm.containerManager.StartContainer(containerId)
		if err != nil {
			fmt.Println("Failed to start container", err)
			return "", err
		}
	}

	// CREATE POD SUCCESS
	pod.Status.Phase = minik8s_pod.PodRunning

	return pod.Metadata.UUID, nil
}

func (rm *RuntimeManager) RemovePod(pod *minik8s_pod.Pod) error {
	// First, remove all containers in the pod
	for _, container := range pod.Spec.Containers {
		_, err := rm.containerManager.RemoveContainer(container.Id)
		if err != nil {
			return err
		}
	}

	// Remove the pause container
	_, err := rm.containerManager.RemoveContainer("pause-" + pod.Metadata.UUID)
	if err != nil {
		return err
	}

	return nil
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

	PodPortBindings := nat.PortMap{}
	PodExposedPorts := nat.PortSet{}

	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.ContainerPort == 0 {
				// If the container port is not set, skip it
				continue
			}
			if port.HostIP == "" {
				port.HostIP = "127.0.0.1"
			}
			if port.Protocol == "" {
				port.Protocol = minik8s_container.ProtocolTCP
			}
			if port.HostPort == 0 {
				// Check if the ContainerPort is already used
				// If not in use, assign it to HostPort.
				// Otherwise, assign the first available port to HostPort.
				containerport_available := nettools.CheckPortAvailability(int(port.ContainerPort))
				if containerport_available {
					port.HostPort = port.ContainerPort
				} else {
					available_port, err := nettools.GetAvailablePort()
					if err != nil {
						fmt.Println("Failed to get available port")
						return "", err
					}
					port.HostPort = (int32)(available_port)
				}
			}

			// Add port binding
			BindingKey, err := nat.NewPort((string)(port.Protocol), fmt.Sprintf("%d", port.ContainerPort))
			if err != nil {
				fmt.Println("Failed to create port binding")
			}

			// Check if the port is already in the map
			if _, ok := PodPortBindings[BindingKey]; ok {
				fmt.Println("Port binding already exists")
			} else {
				// Bind the port
				PodPortBindings[BindingKey] = []nat.PortBinding{
					{
						HostIP:   port.HostIP,
						HostPort: fmt.Sprintf("%d", port.HostPort),
					},
				}
			}

			// Add exposed port
			PodExposedPorts[BindingKey] = struct{}{}
		}
	}

	pauseContainerId, err := rm.containerManager.CreateContainer(pauseName, &minik8s_container.CreateContainerConfig{
		Image:        "registry.aliyuncs.com/google_containers/pause:3.9",
		IpcMode:      "shareable",
		PortBindings: PodPortBindings,
		ExposedPorts: PodExposedPorts,
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

	// Get the pause container's IP address and set it as the pod's IP address
	pauseContainerIp, err := rm.containerManager.GetContainerIpAddress(pauseContainerId)
	if err != nil {
		println("Failed to get pause container's IP address")
		return "", err
	}
	pod.Status.PodIP = pauseContainerIp

	return pauseContainerId, nil
}

func (rm *RuntimeManager) GetNodeStatus() (minik8s_node.NodeStatus, error) {
	hostName, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get hostname")
		return minik8s_node.NodeStatus{}, err
	}

	nodeIp := nettools.KubeletDefaultIP()

	nodeConditions := []string{
		minik8s_node.NODE_Ready,
	}

	cpuPercent := nodestatusutils.GetNodeCpuPercent()
	memPercent := nodestatusutils.GetNodeMemPercent()

	updateTime := time.Now().String()

	// TODO: Number of pods should be filled in the kubelet, the length of map in podManager

	return minik8s_node.NodeStatus{
		Hostname:   hostName,
		Ip:         nodeIp,
		Condition:  nodeConditions,
		CpuPercent: cpuPercent,
		MemPercent: memPercent,
		NumPods:    0,
		UpdateTime: updateTime,
	}, nil

}

// func (rm *RuntimeManager) RuncAdvicorContainer() {
// 	containerConfig := &minik8s_container.CreateContainerConfig{
// 		Image: "gcr.nju.edu.cn/cadvisor/cadvisor:v0.49.1",
// 	}
// }
