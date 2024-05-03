package handler

import (
	"minik8s/pkg/apiserver/process"
	"minik8s/pkg/ty"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin server的 /api/v1/addReplicaSet对应的方法
func AddReplicaSet(c *gin.Context) {
	var desc ty.ReplicaSetDesc
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	process.AddReplicaSet(&desc)
}

func RemoveReplicaSet(c *gin.Context) {
	var namespace string
	var name string

	err := process.RemoveReplicaSet(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func ModifyReplicaSet(c *gin.Context) {
	var namespace string
	var name string

	err := process.ModifyReplicaSet(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func GetReplicaSet(c *gin.Context) {
	var namespace string
	var name string

	result, err := process.GetReplicaSet(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
