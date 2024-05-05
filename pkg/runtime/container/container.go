package container

import (
	"context"
	"encoding/json"

	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/docker"
	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
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

func (cm *ContainerManager) CreateContainer(name string, config *container.Config) (string, error) {
	// Create a new container

	// First, check image exists otherwise pull image to local
	imageManager := &image.ImageManager{}
	_, err := imageManager.PullImage(config.Image)
	if err != nil {
		return "", err
	}

	response, err := docker.DockerClient.ContainerCreate(
		context.Background(),
		config,
		nil,
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
