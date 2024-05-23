package image

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/MiniK8s-SE3356/minik8s/pkg/runtime/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

const (
	PauseContainerImage    = "registry.aliyuncs.com/google_containers/pause:3.9"
	cAdvisorContainerImage = "gcr.nju.edu.cn/cadvisor/cadvisor:v0.49.1"
)

type ImageManager struct{}

func (im *ImageManager) PullImage(imageName string) (string, error) {
	// First, check if the image exists locally
	images, err := im.ListImagesWithRef(imageName)
	if err != nil {
		return "", err
	}
	if len(images) > 0 {
		// If the image exists locally, return the ID
		return images[0].ID, nil
	}

	// Pull an image
	reader, err := docker.DockerClient.ImagePull(
		context.Background(),
		imageName,
		image.PullOptions{},
	)

	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Display the pull progress
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(reader, os.Stderr, termFd, isTerm, nil)
	fmt.Println("Image pulled successfully")

	// Check local images again
	images, err = im.ListImagesWithRef(imageName)
	if err != nil {
		return "", err
	}
	if len(images) > 0 {
		return images[0].ID, nil
	}

	return "", errors.New("failed to pull image")
}

func (im *ImageManager) ListAllImages() ([]image.Summary, error) {
	images, err := docker.DockerClient.ImageList(
		context.Background(),
		image.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (im *ImageManager) ListImagesWithRef(ref string) ([]image.Summary, error) {
	filter := filters.NewArgs()
	filter.Add("reference", im.ImageRefSwitch(ref))

	images, err := docker.DockerClient.ImageList(
		context.Background(),
		image.ListOptions{
			Filters: filter,
		},
	)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// When you pull an image from 'docker.io', the local image reference will be only '<image_name>:<tag>'.
// When you pull an image from a private registry, the local image reference will be '<registry>/<image_name>:<tag>'.
func (im *ImageManager) ImageRefSwitch(ref string) string {
	// If the image reference is not from 'docker.io', return the original reference
	if !strings.HasPrefix(ref, "docker.io/") {
		return ref
	}

	// Split the image reference by '/'
	parts := strings.Split(ref, "/")

	if len(parts) > 1 {
		// If the image reference is from 'docker.io', return the image name with the tag
		return parts[len(parts)-1]
	}

	return ref
}
