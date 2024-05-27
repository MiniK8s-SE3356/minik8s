package cmdline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/server"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/spf13/cobra"
)

func CreateCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Usage()
		return
	}

	// function名称 zip路径
	functionName := args[0]
	zipPath := args[1]

	// 读取zip文件内容
	zipContent, err := os.ReadFile(zipPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var desc struct {
		FunctionName string `json:"functionName"`
		ZipContent   string `json:"zipContent"`
	}

	desc.FunctionName = functionName
	desc.ZipContent = string(zipContent)
	jsonData, _ := json.Marshal(desc)
	result, err := httpRequest.PostRequest(server.RootURL+server.CreateFunction, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result is", result)

	// fmt.Println("result is ", result)
}
