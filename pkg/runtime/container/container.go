package container

import (
	"context"
	"encoding/json"

	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/docker"
	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	minik8sTypes "github.com/MiniK8s-SE3356/minik8s/pkg/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

func (cm *ContainerManager) CreateContainer(name string, config *minik8sTypes.CreateContainerConfig) (string, error) {
	// Create a new container

	// First, check image exists otherwise pull image to local
	imageManager := &image.ImageManager{}
	_, err := imageManager.PullImage(config.Image)
	if err != nil {
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
