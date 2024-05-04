package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/ty"

	"github.com/gin-gonic/gin"
)

func AddPod(c *gin.Context) {
	var desc ty.PodDesc
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddPod(&desc)
	if err != nil {
		fmt.Println("error in process.AddPod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func RemovePod(c *gin.Context) {

}

func ModifyPod(c *gin.Context) {

}

func GetPod(c *gin.Context) {

}

func GetPods(c *gin.Context) {

}
