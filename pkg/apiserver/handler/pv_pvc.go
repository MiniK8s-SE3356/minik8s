package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/gin-gonic/gin"
)

func GetAllPersistVolume(c *gin.Context) {

}

func UpdatePersistVolume(c *gin.Context) {

}

func AddPV(c *gin.Context) {
	requestMsg := httpobject.HTTPRequest_AddPV{}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// FIXME: Namespace只能为default,这明显需要修改
	result, err := process.AddPV(process.DefaultNamespace, &requestMsg.Pv)
	if err != nil {
		fmt.Println("error in process.AddPod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func AddPVC(c *gin.Context) {

}
