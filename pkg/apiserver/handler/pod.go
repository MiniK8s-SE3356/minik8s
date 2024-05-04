package handler

import (
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/ty"

	"github.com/gin-gonic/gin"
)

func AddPod(c *gin.Context) {
	var desc ty.PodDesc
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	process.AddPod(&desc)
}

func RemovePod(c *gin.Context) {

}

func ModifyPod(c *gin.Context) {

}

func GetPod(c *gin.Context) {

}
