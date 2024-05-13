package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/gin-gonic/gin"
)

func GetAllEndpoint(c *gin.Context) {
	result, err := process.GetAllEndpoint()
	if err != nil {
		fmt.Println("error in process.GetAllEndpoint")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func UpdateEndpointBatch(c *gin.Context) {
	var desc map[string]service.EndPoint
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := process.UpdateEndpointBatch(desc)
	if err != nil {
		fmt.Println("error in process.UpdateEndpointBatch")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
