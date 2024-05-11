package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/ty"

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
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace, ok := param["namespace"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'namespace'"})
		return
	}

	name, ok := param["name"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'name'"})
		return
	}

	err := process.RemoveReplicaSet(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// func ModifyReplicaSet(c *gin.Context) {
// 	var namespace string
// 	var name string

// 	err := process.ModifyReplicaSet(namespace, name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}
// }

func GetReplicaSet(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace := c.Query("namespace")

	name := c.Query("name")

	var result string
	var err error
	// 四种情况
	// 1. namespace name均为空 获取全部replicaset
	if namespace == "" && name == "" {
		result, err = process.GetAllReplicaSets()
	}
	// 2. namespace为空，name不为空 获取Default下的replicaset
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.GetReplicaSet(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部replicaset
	if namespace != "" && name == "" {
		result, err = process.GetReplicaSets(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		result, err = process.GetReplicaSet(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.GetReplicaSet ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func DescribeReplicaSet(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	namespace := c.Query("namespace")

	name := c.Query("name")

	var result string
	var err error
	// 四种情况
	// 1. namespace name均为空 获取全部replicaset
	if namespace == "" && name == "" {
		result, err = process.DescribeAllReplicaSets()
	}
	// 2. namespace为空，name不为空 获取Default下的replicaset
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.DescribeReplicaSet(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部replicaset
	if namespace != "" && name == "" {
		result, err = process.DescribeReplicaSets(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		result, err = process.DescribeReplicaSet(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.DescribeReplicaSet ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
