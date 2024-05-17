package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/gin-gonic/gin"
)

// POST 参数类型ServiceDesc
func AddService(c *gin.Context) {
	var requestMsg struct {
		Namespace string
		Desc      map[string]interface{}
	}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddService(requestMsg.Namespace, &requestMsg.Desc)
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

func GetFilteredService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "hello"})
}

func GetAllService(c *gin.Context) {
	result, err := process.GetAllService()

	if err != nil {
		fmt.Println("error in process.GetAllService ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func GetService(c *gin.Context) {
	// var param map[string]string
	// if err := c.ShouldBindJSON(&param); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	namespace := c.Query("namespace")
	name := c.Query("name")

	var result map[string][]interface{}
	var err error
	// 四种情况
	// 1. namespace name均为空 获取全部service
	if namespace == "" && name == "" {
		result, err = process.GetAllService()
	}
	// 2. namespace为空，name不为空 获取Default下的service
	if namespace == "" && name != "" {
		namespace = "Default"
		// result, err = process.GetService(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部service
	if namespace != "" && name == "" {
		// result, err = process.GetServices(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		// result, err = process.GetService(namespace, name)
	}

	if err != nil {
		fmt.Println("error in handler.GetService ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func UpdateService(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var arr []string
	for k, v := range param {
		_, err := process.UpdateService("Default", k, v)
		if err != nil {
			fmt.Println("failed to update", k)
		} else {
			arr = append(arr, k)
		}
	}

	// if err != nil {
	// 	fmt.Println("error in process.UpdateService ", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// }

	c.JSON(http.StatusOK, gin.H{"result": arr})
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
		fmt.Println("error in handler.DescribeService ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}
