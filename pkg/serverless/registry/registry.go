package registry

import (
	"fmt"

	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
)

// Init a registry container and return the registry container id
func RegistryInit() (string, error) {
	containerManager := &minik8s_container.ContainerManager{}
	registryContainer, err := containerManager.GetContainerByName(REGISTRY_NAME)
	if err == nil {
		// The registry container already exists
		if registryContainer.State == "running" {
			// The registry container is running
			fmt.Println("The registry container is already running")
			return registryContainer.ID, nil
		} else {
			registryContainerId, err := containerManager.RestartContainer(registryContainer.ID)
			if err != nil {
				fmt.Println("Error restarting the registry container")
				return "", err
			}
			fmt.Println("The registry container has been restarted")
			return registryContainerId, nil
		}
	}

	// The registry container does not exist
	registryContainerId, err := containerManager.RunDefaultRegistryContainer()
	if err != nil {
		fmt.Println("Error running the registry container")
		return "", err
	}

	return registryContainerId, nil
}
