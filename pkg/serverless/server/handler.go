package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/etcdclient"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/workflow"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/idgenerate"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
	"github.com/gin-gonic/gin"
)

var EtcdCli *etcdclient.EtcdClient
var Mq *message.MQConnection

const functionPrefix = "/minik8s/function/"

func triggerServerlessFunction(c *gin.Context) {
	var desc struct {
		FunctionName string `json:"functionName"`
		Params       string `json:"params"`
	}
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// result, err :=
	result := ""
	// if err != nil {
	// 	fmt.Println("error in triggerServerlessWorkflow ", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// }

	c.JSON(http.StatusOK, result)
}

func triggerServerlessWorkflow(c *gin.Context) {
	var desc struct {
		Workflow workflow.Workflow `json:"workflow"`
		MqName   string            `json:"mqName"`
	}
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// result, err :=
	result := ""
	// if err != nil {
	// 	fmt.Println("error in triggerServerlessWorkflow ", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// }

	c.JSON(http.StatusOK, result)
}

func createFunction(c *gin.Context) {
	var desc struct {
		FunctionName string `json:"functionName"`
		ZipContent   string `json:"zipContent"`
	}
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var f function.Function
	f.APIVersion = "v1"
	f.Kind = "Function"
	uuid, _ := idgenerate.GenerateID()
	f.Metadata.UUID = uuid
	f.Metadata.Name = desc.FunctionName
	f.Metadata.Namespace = "Default"
	f.Spec.ImageName = ""
	result, err := createFunctionProcess(f)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func createFunctionProcess(f function.Function) (string, error) {
	result := ""
	existed, err := EtcdCli.Exist(functionPrefix + f.Metadata.Name)
	if err != nil {
		return result, err
	}

	if existed {
		return result, nil
	}

	value, err := json.Marshal(f)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return result, err
	}

	// 然后存入etcd
	err = EtcdCli.Put(functionPrefix+f.Metadata.Name, string(value))
	if err != nil {
		fmt.Println("failed to write to etcd ", err.Error())
		return result, err
	}

	return "create function success", nil
}
