package cmdline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	gpu_types "github.com/MiniK8s-SE3356/minik8s/pkg/gpu/types"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var SubmitFuncTable = map[string]func(args []string) error{
	"GPUJob": SubmitGPUJob,
}

func SubmitCmdHandler(cmd *cobra.Command, args []string) {
	kind := args[0]
	submitFunc, ok := SubmitFuncTable[kind]
	if !ok {
		fmt.Println("kind not supported")
		return
	}

	err := submitFunc(args[1:])
	if err != nil {
		fmt.Println("error in SubmitCmdHandler ", err.Error())
		return
	}
}

func SubmitGPUJob(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("args number is not correct")
	}
	yamlPath := args[0]
	jobZipPath := args[1]

	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		fmt.Println("failed to read yaml file")
		return err
	}

	var gpuSlurmJobDesc gpu_types.SlurmJob
	err = yaml.Unmarshal(yamlContent, &gpuSlurmJobDesc)
	if err != nil {
		fmt.Println("failed to unmarshal yaml file")
		return err
	}

	//!debug//
	debugJson, _ := json.Marshal(gpuSlurmJobDesc)
	fmt.Println("gpuSlurmJobDesc is", string(debugJson))
	//!debug//

	uuid, err := idgenerate.GenerateID()
	if err != nil {
		fmt.Println("failed to generate id")
		return err
	}
	gpuSlurmJobDesc.Metadata.UUID = uuid

	zipContent, err := os.ReadFile(jobZipPath)
	if err != nil {
		fmt.Println("failed to read zip file")
		return err
	}

	// Create a job request
	var jobRequest struct {
		JobDesc    gpu_types.SlurmJob `json:"jobDesc"`
		ZipContent []byte             `json:"zipContent"`
	}
	jobRequest.JobDesc = gpuSlurmJobDesc
	jobRequest.ZipContent = zipContent
	jsonData, _ := json.Marshal(jobRequest)
	result, err := httpRequest.PostRequest(
		GPUCtlRootURL+url.SubmitGPUJob,
		jsonData,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("result is", result)
	return nil
}
