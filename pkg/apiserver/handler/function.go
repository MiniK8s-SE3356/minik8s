package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/gin-gonic/gin"
)

func GetAllServerlessFunction(c *gin.Context) {
	functionNames, err := process.GetAllServerlessFunction()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, functionNames)
}
