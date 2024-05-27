package registry_test

import (
	// "encoding/json"
	"testing"
	// "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	// runtime_image "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	// "github.com/MiniK8s-SE3356/minik8s/pkg/serverless/registry"
	// "github.com/docker/docker/api/types/image"
)

func TestMain(m *testing.M) {
	// registryContainerId, err := registry.RegistryInit()
	// if err != nil {
	// 	println("Error initializing the registry container")
	// 	panic(err)
	// }
	// imageManager := &runtime_image.ImageManager{}
	// containerManager := &container.ContainerManager{}

	// _, err = imageManager.PullImage("ubuntu:latest")
	// if err != nil {
	// 	println("Error pulling the image")
	// 	panic(err)
	// }
	// err = imageManager.TagImage("ubuntu:latest", "localhost:5000/ubuntu:latest")
	// if err != nil {
	// 	println("Error tagging the image")
	// 	panic(err)
	// }

	// authStr, _ := json.Marshal(registry.REGISTRY_AUTH_CONFIG)
	// options := image.PushOptions{
	// 	All:          false,
	// 	RegistryAuth: string(authStr),
	// }
	// err = imageManager.PushImageWithOptions("localhost:5000/ubuntu:latest", options)
	// if err != nil {
	// 	println("Error pushing the image")
	// 	panic(err)
	// }

	// _, err = containerManager.ForceRemoveContainer(registryContainerId)
	// if err != nil {
	// 	println("Error removing the registry container")
	// 	panic(err)
	// }
}
