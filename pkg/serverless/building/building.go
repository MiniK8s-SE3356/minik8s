package building

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	minik8s_docker "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/docker"
	runtime_image "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/registry"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
	minik8s_zip "github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

const (
	BASE_IMAGE = "levixubbbb/serverless-base-image:latest"
)

func BuildServerlessFunctionImage(function *function.Function) error {
	// Create a new directory for building the function image
	err := os.Mkdir(function.Metadata.UUID, 0777)
	if err != nil {
		fmt.Println("Error creating directory for function image building")
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory")
		return err
	}

	err = os.WriteFile(
		path.Join(pwd, function.Metadata.UUID, "serverless.zip"),
		function.Spec.FileContent,
		os.ModePerm,
	)
	if err != nil {
		fmt.Println("Error writing function content to file")
		return err
	}

	err = minik8s_zip.DecompressZipFile(
		path.Join(pwd, function.Metadata.UUID, "serverless.zip"),
		path.Join(pwd, function.Metadata.UUID),
	)
	if err != nil {
		fmt.Println("Error decompressing function zip file")
		return err
	}

	// Create a Dockerfile for building the function image
	dockerfile, err := os.Create(path.Join(pwd, function.Metadata.UUID, "Dockerfile"))
	if err != nil {
		fmt.Println("Error creating Dockerfile for function image building")
		return err
	}

	dockerfile.WriteString("FROM " + BASE_IMAGE + "\n")
	// srcFile := path.Join(function.Metadata.UUID, function.Spec.FilePath, "*")
	srcFile := path.Join(function.Spec.FilePath, "*")
	dstFile := "/app/"
	dockerfile.WriteString("COPY " + srcFile + " " + dstFile + "\n")
	dockerfile.WriteString("RUN pip install --no-cache-dir -r requirements.txt\n")
	dockerfile.WriteString("EXPOSE 5000\n")

	dockerfile.Close()

	// Build the function image
	dockerClient := minik8s_docker.NewDockerClient()

	// Docker need a .tar file to build the image
	tarFilename := "function-" + function.Metadata.UUID + ".tar"
	err = minik8s_zip.CompressTarFile(
		path.Join(pwd, function.Metadata.UUID),
		path.Join(pwd, tarFilename),
	)
	if err != nil {
		fmt.Println("Error compressing function directory to tar file")
		return err
	}

	tarFile, err := os.Open(path.Join(pwd, tarFilename))
	if err != nil {
		fmt.Println("Error opening tar file")
		return err
	}
	defer tarFile.Close()

	options := types.ImageBuildOptions{
		Tags:        []string{registry.REGISTRY_IP + ":" + registry.REGISTRY_PORT + "/func-" + function.Metadata.UUID + ":latest"},
		Dockerfile:  "Dockerfile", //path.Join(function.Metadata.UUID, "Dockerfile"),
		Remove:      true,
		ForceRemove: true,
		Context:     tarFile,
	}

	response, err := dockerClient.ImageBuild(
		context.Background(),
		tarFile,
		options,
	)
	if err != nil {
		fmt.Println("Error building function image")
		return err
	}
	defer response.Body.Close()
	fmt.Println("Function image built successfully")

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		fmt.Println("Error copying response body")
		return err
	}

	// TODO: Push the image to a registry
	imageManager := runtime_image.ImageManager{}
	authStr, _ := json.Marshal(registry.REGISTRY_AUTH_CONFIG)

	err = imageManager.PushImageWithOptions(
		registry.REGISTRY_IP+":"+registry.REGISTRY_PORT+"/func-"+function.Metadata.UUID+":latest",
		image.PushOptions{
			All:          false,
			RegistryAuth: string(authStr),
		},
	)

	if err != nil {
		fmt.Println("Error pushing function image to registry")
		return err
	}

	// Remove the temporary directory and tar file
	err = os.RemoveAll(path.Join(pwd, function.Metadata.UUID))
	if err != nil {
		fmt.Println("Error removing temporary directory")
		return err
	}

	err = os.RemoveAll(path.Join(pwd, tarFilename))
	if err != nil {
		fmt.Println("Error removing tar file")
		return err
	}

	err = imageManager.RemoveImageByName(registry.REGISTRY_IP + ":" + registry.REGISTRY_PORT + "/func-" + function.Metadata.UUID + ":latest")
	if err != nil {
		fmt.Println("Error removing function image")
		return err
	}

	return nil
}
