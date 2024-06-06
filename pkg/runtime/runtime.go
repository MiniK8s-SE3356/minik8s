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

var NodeRuntimeMangaer = NewRuntimeManager()

type RuntimeManager struct {
	containerManager *runtime_container.ContainerManager
	imageManager     *runtime_image.ImageManager
}

func NewRuntimeManager() RuntimeManager {
	return RuntimeManager{
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

		//// Rename the container to pod-uuid-container-name
		// container.Name = pod.Metadata.UUID + "-" + container.Name

		// parse environment variables
		env := []string{}
		for _, envVar := range container.Env {
			env = append(env, fmt.Sprintf("%s=%s", envVar.Name, envVar.Value))
		}

		// parse volume mounts. We should mount the pod's volumes to the container
		volumeBinds := rm.GetVolumeBinds(&pod.Spec.Volumes, &container.VolumeMounts)
		for _, volumeBind := range volumeBinds {
			fmt.Println("Volume bind: ", volumeBind)
		}

		//! Used for DNS resolution, bind the '/etc/hosts' file to the container
		volumeBinds = append(volumeBinds, "/etc/hosts:/etc/hosts:ro")

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
			Cmd:         container.Command,
		}
		containerId, err := rm.containerManager.CreateContainer(container.Name, containerConfig)
		if err != nil {
			fmt.Println("Failed to create container", err)
			return "", err
		}
		// Change pod's container id and name
		pod.Spec.Containers[i].Id = containerId
		// pod.Spec.Containers[i].Name = container.Name
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
	_, err := rm.imageManager.PullImage(runtime_image.PauseContainerImage)
	if err != nil {
		return "", err
	}

	uuid := pod.Metadata.UUID
	// A pod's pause container is named as "pause-<pod-uuid>"
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
			//! If user don't set HostIP, "0.0.0.0" will be used
			// if port.HostIP == "" {
			// 	port.HostIP = "127.0.0.1"
			// }
			if port.Protocol == "" {
				port.Protocol = minik8s_container.ProtocolTCP
			}
			// if port.HostPort == 0 {
			// Check if the ContainerPort is already used
			// If not in use, assign it to HostPort.
			// Otherwise, assign the first available port to HostPort.
			//! If user don't set HostPort, we don't expose the port to the host
			// containerport_available := nettools.CheckPortAvailability(int(port.ContainerPort))
			// if containerport_available {
			// 	port.HostPort = port.ContainerPort
			// } else {
			// 	available_port, err := nettools.GetAvailablePort()
			// 	if err != nil {
			// 		fmt.Println("Failed to get available port")
			// 		return "", err
			// 	}
			// 	port.HostPort = (int32)(available_port)
			// }
			// }

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
				var hostPort string
				if port.HostPort == 0 {
					hostPort = ""
				} else {
					hostPort = fmt.Sprintf("%d", port.HostPort)
				}
				PodPortBindings[BindingKey] = []nat.PortBinding{
					{
						HostIP:   port.HostIP,
						HostPort: hostPort,
					},
				}
			}

			// Add exposed port
			PodExposedPorts[BindingKey] = struct{}{}
		}
	}

	//! Used for DNS resolution, bind the '/etc/hosts' file to the container
	volumes := []string{"/etc/hosts:/etc/hosts:ro"}

	pauseContainerId, err := rm.containerManager.CreateContainer(pauseName, &minik8s_container.CreateContainerConfig{
		Image:        runtime_image.PauseContainerImage,
		IpcMode:      "shareable",
		PortBindings: PodPortBindings,
		ExposedPorts: PodExposedPorts,
		Binds:        volumes,
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

	updateTime := time.Now()

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

func (rm *RuntimeManager) RestartPod(pod *minik8s_pod.Pod) (string, error) {

	fmt.Println("Restarting pod for ", pod.Metadata.Name, " with UUID ", pod.Metadata.UUID)
	// Check this pod's pause container exists or not
	pauseName := "pause-" + pod.Metadata.UUID
	pauseContainer, err := rm.containerManager.GetContainerByName(pauseName)
	var pauseContainerId string
	if err != nil {
		// This means the pause container does not exist, we have to create it again
		pauseContainerId, err = rm.CreateAndStartPauseContainer(pod)
		if err != nil {
			return "", err
		}
	} else {
		pauseContainerId, err = rm.containerManager.RestartContainer(pauseContainer.ID)
		if err != nil {
			return "", err
		}
	}

	// Pause container's IP address maybe changed, we have to update the pod's IP address
	pod.Status.PodIP, err = rm.containerManager.GetContainerIpAddress(pauseContainerId)
	if err != nil {
		return "", err
	}

	// Restart all containers in the pod
	for i, container := range pod.Spec.Containers {
		// Check this container exists or not
		realContainer, err := rm.containerManager.GetContainerByName(container.Name)
		if err != nil {
			env := []string{}
			for _, envVar := range container.Env {
				env = append(env, fmt.Sprintf("%s=%s", envVar.Name, envVar.Value))
			}

			// TODO: parse command

			volumeBinds := rm.GetVolumeBinds(&pod.Spec.Volumes, &container.VolumeMounts)
			for _, volumeBind := range volumeBinds {
				fmt.Println("Volume bind: ", volumeBind)
			}

			containerConfig := &minik8s_container.CreateContainerConfig{
				Image:       container.Image,
				NetworkMode: "container:" + pauseContainerId,
				IpcMode:     "container:" + pauseContainerId,
				PidMode:     "container:" + pauseContainerId,
				Env:         env,
				Cmd:         container.Command,
				Binds:       volumeBinds,
				CPU:         container.Resources.Limits.CPU,
				Memory:      container.Resources.Limits.Memory,
			}
			containerId, err := rm.containerManager.CreateContainer(container.Name, containerConfig)
			if err != nil {
				fmt.Println("Failed to create container", err)
				return "", err
			}
			pod.Spec.Containers[i].Id = containerId
			// Start the container
			_, err = rm.containerManager.StartContainer(containerId)
			if err != nil {
				fmt.Println("Failed to start container", err)
				return "", err
			}
		} else {
			_, err = rm.containerManager.RestartContainer(realContainer.ID)
			if err != nil {
				fmt.Println("Failed to restart container: ", err)
				return "", err
			}
		}
	}

	pod.Status.Phase = minik8s_pod.PodRunning

	return pod.Metadata.UUID, nil
}

func (rm *RuntimeManager) RuncAdvicorContainer() {
	// Check the cAdvisor container exits or not
	cAdvisorContainerName := "cadvisor"
	cAdvisorContainer, err := rm.containerManager.GetContainerByName(cAdvisorContainerName)
	if err == nil {
		if cAdvisorContainer.State == "running" {
			// cAdvisor has already been running
			fmt.Println("cAdvisor container has already been running! ID: ", cAdvisorContainer.ID)
			return
		} else {
			cAdvisorId, err := rm.containerManager.RestartContainer(cAdvisorContainer.ID)
			if err == nil {
				fmt.Println("cAdvisor container restart success! ID: ", cAdvisorId)
				return
			}
		}
	}

	cAdvisorId, err := rm.containerManager.RunDefaultcAdvisorContainer()
	if err != nil {
		fmt.Println("RunDefaultcAdvisorContainer failed")
		return
	}
	fmt.Println("cAdvisor container start success! ID: ", cAdvisorId)
}
