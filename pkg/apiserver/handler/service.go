package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/ty"
	"github.com/gin-gonic/gin"
)

// POST 参数类型ServiceDesc
func AddService(c *gin.Context) {
	var desc ty.ServiceDesc
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddService(&desc)
	if err != nil {
		fmt.Println("error in process.AddService ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// POST 参数类型 {namespace: "xxx", name: "xxx"}
func RemoveService(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace, ok := param["namespace"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'name'"})
		return
	}

	name, ok := param["name"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'name'"})
		return
	}

	result, err := process.RemoveService(namespace, name)
	if err != nil {
		fmt.Println("error in process.RemoveService ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func GetService(c *gin.Context) {
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
	// 1. namespace name均为空 获取全部service
	if namespace == "" && name == "" {
		result, err = process.GetAllServices()
	}
	// 2. namespace为空，name不为空 获取Default下的service
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.GetService(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部service
	if namespace != "" && name == "" {
		result, err = process.GetServices(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		result, err = process.GetService(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.GetService ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func DescribeService(c *gin.Context) {
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
	// 1. namespace name均为空 获取全部service
	if namespace == "" && name == "" {
		result, err = process.DescribeAllServices()
	}
	// 2. namespace为空，name不为空 获取Default下的service
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.DescribeService(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部service
	if namespace != "" && name == "" {
		result, err = process.DescribeServices(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		result, err = process.DescribeService(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.DescribeService ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
