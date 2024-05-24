package building

import (
	"fmt"
	"os"
	"path"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serveless/types/function"
	minik8s_zip "github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
)

func BuildServelessFunctionImage(function *function.Function) error {
	// Create a new directory for building the function image
	err := os.Mkdir(function.Metadata.Name, 0777)
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
		path.Join(pwd, function.Metadata.Name, "serveless.zip"),
		function.Spec.FileContent,
		os.ModePerm,
	)
	if err != nil {
		fmt.Println("Error writing function content to file")
		return err
	}

	err = minik8s_zip.DecompressZipFile(
		path.Join(pwd, function.Metadata.Name, "serveless.zip"),
		path.Join(pwd, function.Metadata.Name),
	)
	if err != nil {
		fmt.Println("Error decompressing function zip file")
		return err
	}

	// Create a Dockerfile for building the function image
	dockerfile, err := os.Create(path.Join(pwd, function.Metadata.Name, "Dockerfile"))
	if err != nil {
		fmt.Println("Error creating Dockerfile for function image building")
		return err
	}

	dockerfile.WriteString("FROM levixubbbb/serveless-base-image:latest\n")
	dockerfile.WriteString("COPY . /app\n")
	dockerfile.WriteString("EXPOSE 5000\n")

	dockerfile.Close()

	// Build the function image

	return nil
}
