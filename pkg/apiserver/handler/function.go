package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
	"github.com/gin-gonic/gin"
)

func GetFunction(c *gin.Context) {
	name := c.Query("name")
	var result []function.Function
	var err error
	if name == "" {
		result, err = process.GetAllFunction()
	} else {
		result, err = process.GetFunction(name)
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetAllServerlessFunction(c *gin.Context) {
	functionNames, err := process.GetAllServerlessFunction()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, functionNames)
}
