package container

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/docker"
	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

type ContainerManager struct{}

func (cm *ContainerManager) GetContainerStats(id string) (*types.StatsJSON, error) {
	// Get a container status
	stats, err := docker.DockerClient.ContainerStats(
		context.Background(),
		id,
		false,
	)

	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	containerStats := &types.StatsJSON{}
	err = json.NewDecoder(stats.Body).Decode(containerStats)

	if err != nil {
		return nil, err
	}

	return containerStats, nil
}

func (cm *ContainerManager) GetContainerInspect(id string) (*types.ContainerJSON, error) {
	// Get a container inspect info
	containerInfo, err := docker.DockerClient.ContainerInspect(
		context.Background(),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &containerInfo, nil
}

func (cm *ContainerManager) GetContainerIpAddress(id string) (string, error) {
	// Get a container IP address
	containerInfo, err := cm.GetContainerInspect(id)

	if err != nil {
		return "", err
	}

	return containerInfo.NetworkSettings.IPAddress, nil
}

func (cm *ContainerManager) ListAllContainers() ([]types.Container, error) {
	// List all containers
	containers, err := docker.DockerClient.ContainerList(
		context.Background(),
		container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(),
		},
	)

	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (cm *ContainerManager) CreateContainer(name string, config *minik8s_container.CreateContainerConfig) (string, error) {
	// Create a new container

	// First, check image exists otherwise pull image to local
	imageManager := &image.ImageManager{}
	_, err := imageManager.PullImage(config.Image)
	if err != nil {
		fmt.Println("Failed to pull image, err:", err)
		return "", err
	}

	response, err := docker.DockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        config.Image,
			Volumes:      config.Volumes,
			WorkingDir:   config.WorkingDir,
			Env:          config.Env,
			Entrypoint:   config.Entrypoint,
			Cmd:          config.Cmd,
			ExposedPorts: config.ExposedPorts,
			Labels:       config.Labels,
		},
		&container.HostConfig{
			Binds:        config.Binds,
			PortBindings: config.PortBindings,
			NetworkMode:  container.NetworkMode(config.NetworkMode),
			IpcMode:      container.IpcMode(config.IpcMode),
			PidMode:      container.PidMode(config.PidMode),
			VolumesFrom:  config.VolumesFrom,
			Resources: container.Resources{
				NanoCPUs: config.CPU,
				Memory:   config.Memory,
			},
		},
		nil,
		nil,
		name,
	)

	if err != nil {
		return "", err
	}

	return response.ID, nil
}

func (cm *ContainerManager) RunDefaultRegistryContainer() (string, error) {
	imageManager := &image.ImageManager{}
	_, err := imageManager.PullImage("registry:2")
	if err != nil {
		fmt.Println("Failed to pull image, err:", err)
		return "", err
	}

	portBindings := nat.PortMap{}
	portBindings["5000/tcp"] = []nat.PortBinding{
		{
			HostPort: "5000",
		},
	}
	exposedPorts := nat.PortSet{}
	exposedPorts["5000/tcp"] = struct{}{}

	response, err := docker.DockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        "registry:2",
			ExposedPorts: exposedPorts,
			Tty:          true,
		},
		&container.HostConfig{
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		},
		nil,
		nil,
		"minik8s-registry",
	)

	if err != nil {
		fmt.Println("Registry create failed")
		return "", err
	}

	err = docker.DockerClient.ContainerStart(
		context.Background(),
		response.ID,
		container.StartOptions{},
	)

	if err != nil {
		fmt.Println("Registry start failed")
		return "", err
	}

	return response.ID, nil
}

func (cm *ContainerManager) RunDefaultcAdvisorContainer() (string, error) {
	imageManager := &image.ImageManager{}
	_, err := imageManager.PullImage("gcr.nju.edu.cn/cadvisor/cadvisor:v0.49.1")
	if err != nil {
		fmt.Println("Failed to pull image, err:", err)
		return "", err
	}
	portBindings := nat.PortMap{}
	portBindings["8080/tcp"] = []nat.PortBinding{
		{
			// HostIP:   "127.0.0.1",

			// We use host's '8090' port to run cAdvisor!
			HostPort: "8090",
		},
	}
	exposedPorts := nat.PortSet{}
	exposedPorts["8080/tcp"] = struct{}{}

	response, err := docker.DockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        "gcr.nju.edu.cn/cadvisor/cadvisor:v0.49.1",
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			Binds: []string{
				"/:/rootfs:ro",
				"/var/run:/var/run:ro",
				"/sys:/sys:ro",
				"/var/lib/docker/:/var/lib/docker:ro",
				"/dev/disk/:/dev/disk:ro",
			},
			PortBindings: portBindings,
			Privileged:   true,
			Mounts: []mount.Mount{
				{
					Source: "/dev/kmsg",
					Target: "/dev/kmsg",
					Type:   mount.TypeBind,
				},
			},
		},
		nil,
		nil,
		"cadvisor",
	)

	if err != nil {
		fmt.Println("cAdvisor create failed")
		return "", err
	}

	err = docker.DockerClient.ContainerStart(
		context.Background(),
		response.ID,
		container.StartOptions{},
	)

	if err != nil {
		fmt.Println("cAdvisor start failed")
		return "", err
	}

	return response.ID, nil
}

func (cm *ContainerManager) RemoveContainer(id string) (string, error) {
	// Remove a container
	containerInfo, err := cm.GetContainerInspect(id)

	if err != nil {
		return "", err
	}

	if (containerInfo.State != nil) && (containerInfo.State.Running) {
		_, err = cm.StopContainer(id)

		if err != nil {
			return "", err
		}
	}

	err = docker.DockerClient.ContainerRemove(
		context.Background(),
		id,
		container.RemoveOptions{},
	)

	if err != nil {
		return "", err
	}

	fmt.Println("Container removed successfully")
	return id, nil
}

func (cm *ContainerManager) ForceRemoveContainer(id string) (string, error) {
	// Remove a container
	containerInfo, err := cm.GetContainerInspect(id)

	if err != nil {
		return "", err
	}

	if (containerInfo.State != nil) && (containerInfo.State.Running) {
		_, err = cm.StopContainer(id)

		if err != nil {
			return "", err
		}
	}

	err = docker.DockerClient.ContainerRemove(
		context.Background(),
		id,
		container.RemoveOptions{
			Force: true,
		},
	)

	if err != nil {
		return "", err
	}

	fmt.Println("Container force removed successfully")
	return id, nil
}

func (cm *ContainerManager) StartContainer(id string) (string, error) {
	// Start a container
	err := docker.DockerClient.ContainerStart(
		context.Background(),
		id,
		container.StartOptions{},
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (cm *ContainerManager) RestartContainer(id string) (string, error) {
	// Restart a container
	err := docker.DockerClient.ContainerRestart(
		context.Background(),
		id,
		// We will wait the container for 10 seconds, then kill it
		container.StopOptions{},
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (cm *ContainerManager) StopContainer(id string) (string, error) {
	// Stop a container
	err := docker.DockerClient.ContainerStop(
		context.Background(),
		id,
		container.StopOptions{
			Signal:  "",
			Timeout: nil,
		},
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

// Get the first container with the name
func (cm *ContainerManager) GetContainerByName(name string) (*types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", name)

	containers, err := docker.DockerClient.ContainerList(
		context.Background(),
		container.ListOptions{
			All:     true,
			Filters: filterArgs,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, container := range containers {
		for _, containerName := range container.Names {
			if containerName == "/"+name {
				return &container, nil
			}
		}
		// if container.Names[0] == "/"+name {
		// 	return &container, nil
		// }
	}

	return nil, fmt.Errorf("container with name %s not found", name)
}
