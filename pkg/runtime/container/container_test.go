package container_test

import (
	"testing"

	runtime_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	"github.com/docker/docker/api/types/container"
)

func TestMain(m *testing.M) {
	containerManager := &runtime_container.ContainerManager{}
	// Test create container
	config := &container.Config{
		Image: "nginx:alpine",
	}
	id, err := containerManager.CreateContainer(
		"test",
		config,
	)
	if err != nil {
		println("Error creating container")
		panic(err)
	}

	// Test start container
	id, err = containerManager.StartContainer(id)
	if err != nil {
		println("Error starting container")
		panic(err)
	}

	// Test stop container
	id, err = containerManager.StopContainer(id)
	if err != nil {
		println("Error stopping container")
		panic(err)
	}

	// Test remove container
	id, err = containerManager.RemoveContainer(id)
	if err != nil {
		println("Error removing container")
		panic(err)
	}
}
