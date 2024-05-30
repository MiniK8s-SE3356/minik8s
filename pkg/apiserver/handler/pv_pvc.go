package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/gin-gonic/gin"
)

func GetAllPersistVolume(c *gin.Context) {
	pvresult, err := process.GetAllPersistVolume()
	if err != nil {
		fmt.Println("error in process.GetAllPersistVolume ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	pvcresult, err := process.GetAllPersistVolumeClaim()
	if err != nil {
		fmt.Println("error in process.GetAllPersistVolumeClaim ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	result := httpobject.HTTPResponse_GetAllPersistVolume{}
	result.Pv = pvresult
	result.Pvc = pvcresult
	c.JSON(http.StatusOK, result)
}

func UpdatePersistVolume(c *gin.Context) {
	requestMsg := httpobject.HTTPRequest_UpdatePersistVolume{}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var arr []string
	for k, v := range requestMsg.Pv {
		vbyte,err:=json.Marshal(v)
		if(err!=nil){
			fmt.Println("UpdatePV marshal error:",err.Error())
		}
		_, err = process.UpdatePersistVolume("Default", k, string(vbyte))
		if err != nil {
			fmt.Println("failed to update pv", k)
		} else {
			arr = append(arr, k)
		}
	}
	for k, v := range requestMsg.Pvc {
		vbyte,err:=json.Marshal(v)
		if(err!=nil){
			fmt.Println("UpdatePVC marshal error:",err.Error())
		}
		_, err = process.UpdatePersistVolumeClaim("Default", k, string(vbyte))
		if err != nil {
			fmt.Println("failed to update pv", k)
		} else {
			arr = append(arr, k)
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": arr})
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
		fmt.Println("error in process.AddPV ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func AddPVC(c *gin.Context) {
	requestMsg := httpobject.HTTPRequest_AddPVC{}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// FIXME: Namespace只能为default,这明显需要修改
	result, err := process.AddPVC(process.DefaultNamespace, &requestMsg.Pvc)
	if err != nil {
		fmt.Println("error in process.AddPVC ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
