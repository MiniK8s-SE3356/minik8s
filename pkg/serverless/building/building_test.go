package building_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	apiobject "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject"
	runtime_container "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/container"
	runtime_image "github.com/MiniK8s-SE3356/minik8s/pkg/runtime/image"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/building"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/registry"
	serverless_function "github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
	minik8s_container "github.com/MiniK8s-SE3356/minik8s/pkg/types/container"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/zip"
	"github.com/google/uuid"
)

type Message struct {
	Params string `json:"params"`
}

type Params struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func TestMain(m *testing.M) {
	registryContainerId, err := registry.RegistryInit()
	if err != nil {
		println("Error initializing the registry container")
		panic(err)
	}

	dirName := "test-function"
	subDirName := "function"
	fileName := "function.py"
	message := "import json\n\ndef function(params):\n\tparams = json.loads(params)\n\tx = params['x']\n\ty = params['y']\n\tx = x + y\n\tresp = {\n\t\t'sum': x\n\t}\n\treturn json.dumps(resp)\n#################"
	err = os.Mkdir(dirName, os.ModePerm)
	if err != nil {
		println("Error creating directory")
		panic(err)
	}

	err = os.Mkdir(filepath.Join(dirName, subDirName), os.ModePerm)
	if err != nil {
		println("Error creating sub directory")
		panic(err)
	}

	file, err := os.Create(filepath.Join(dirName, subDirName, fileName))
	if err != nil {
		println("Error creating file")
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		println("Error writing to file")
		panic(err)
	}

	err = zip.CompressZipFile(dirName, "serverless.zip")
	if err != nil {
		println("Error compressing zip file")
		panic(err)
	}

	content, err := os.ReadFile("serverless.zip")
	if err != nil {
		println("Error reading file")
		panic(err)
	}

	function := serverless_function.Function{
		Basic: apiobject.Basic{
			APIVersion: "v1",
			Kind:       "Function",
			Metadata: apiobject.Metadata{
				Name: "test-function",
				UUID: uuid.New().String(),
			},
		},
		Spec: serverless_function.FunctionSpec{
			FileContent: content,
			FilePath:    "function",
		},
	}

	err = os.RemoveAll("serverless.zip")
	if err != nil {
		println("Error removing zip file")
		panic(err)
	}

	err = os.RemoveAll(dirName)
	if err != nil {
		println("Error removing directory")
		panic(err)
	}

	err = building.BuildServerlessFunctionImage(&function)
	if err != nil {
		println("Error building function image")
		panic(err)
	}

	containerManager := runtime_container.ContainerManager{}
	functionContainerId, err := containerManager.CreateContainer("test-function", &minik8s_container.CreateContainerConfig{
		Image: registry.REGISTRY_IP + ":" + registry.REGISTRY_PORT + "/func-" + function.Metadata.UUID + ":latest",
	})
	if err != nil {
		println("Error creating function container")
		panic(err)
	}

	_, err = containerManager.StartContainer(functionContainerId)
	if err != nil {
		println("Error starting function container")
		panic(err)
	}

	functionContainerIP, err := containerManager.GetContainerIpAddress(functionContainerId)
	if err != nil {
		println("Error getting function container IP address")
		panic(err)
	}

	println("Function container IP address: " + functionContainerIP)

	// Wait for the serverless function to start
	println("Waiting for the serverless function to start...")
	time.Sleep(10 * time.Second)

	// Try to call the serverless function
	params := Params{
		X: 114514,
		Y: 42,
	}
	paramsJson, _ := json.Marshal(params)

	requestMessage := Message{
		Params: string(paramsJson),
	}
	fmt.Println("Request message: ", string(paramsJson))
	messageContent, _ := json.Marshal(requestMessage)

	response, err := httpRequest.PostRequest(
		"http://"+functionContainerIP+":5000/api/v1/callfunc",
		messageContent,
	)
	if err != nil {
		println("Error sending request")
		panic(err)
	}
	println(response)

	// Remove all container
	_, err = containerManager.ForceRemoveContainer(functionContainerId)
	if err != nil {
		println("Error removing function container")
		panic(err)
	}

	_, err = containerManager.ForceRemoveContainer(registryContainerId)
	if err != nil {
		println("Error removing registry container")
		panic(err)
	}

	// Remove Image
	imageManager := runtime_image.ImageManager{}
	err = imageManager.RemoveImageByName(registry.REGISTRY_IP + ":" + registry.REGISTRY_PORT + "/func-" + function.Metadata.UUID + ":latest")
	if err != nil {
		println("Error removing function image")
		panic(err)
	}
}
