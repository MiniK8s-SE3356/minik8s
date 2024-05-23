package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"

	"github.com/gin-gonic/gin"
)

// gin server的 /api/v1/addHPA对应的方法
func AddHPA(c *gin.Context) {
	var requestMsg struct {
		Namespace string
		Desc      yaml.HPADesc
	}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	process.AddHPA(requestMsg.Namespace, &requestMsg.Desc)
}

func RemoveHPA(c *gin.Context) {
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

	result, err := process.RemoveHPA(namespace, name)
	if err != nil {
		fmt.Println("error in process.RemoveHPA ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// func ModifyHPA(c *gin.Context) {
// 	var namespace string
// 	var name string

// 	err := process.ModifyHPA(namespace, name)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}
// }

func GetHPA(c *gin.Context) {
	// namespace := c.Query("namespace")
	namespace := process.DefaultNamespace
	name := c.Query("name")

	// 四种情况
	// 1. namespace name均为空 获取全部hpa
	if namespace == "" && name == "" {
		result, err := process.GetAllHPAs()

		if err != nil {
			fmt.Println("error in process.GetHPA ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 2. namespace为空，name不为空 获取Default下的hpa
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err := process.GetHPA(namespace, name)

		if err != nil {
			fmt.Println("error in process.GetHPA ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部hpa
	if namespace != "" && name == "" {
		result, err := process.GetHPAs(namespace)

		if err != nil {
			fmt.Println("error in process.GetHPA ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name != "" {
		result, err := process.GetHPA(namespace, name)

		if err != nil {
			fmt.Println("error in process.GetHPA ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}

}

// func DescribeHPA(c *gin.Context) {
// 	var param map[string]string
// 	if err := c.ShouldBindJSON(&param); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	namespace := c.Query("namespace")

// 	name := c.Query("name")

// 	var result string
// 	var err error
// 	// 四种情况
// 	// 1. namespace name均为空 获取全部hpa
// 	if namespace == "" && name == "" {
// 		result, err = process.DescribeAllHPAs()
// 	}
// 	// 2. namespace为空，name不为空 获取Default下的hpa
// 	if namespace == "" && name != "" {
// 		namespace = "Default"
// 		result, err = process.DescribeHPA(namespace, name)
// 	}
// 	// 3. namespace不为空，name为空 获取namespace下的全部hpa
// 	if namespace != "" && name == "" {
// 		result, err = process.DescribeHPAs(namespace)
// 	}
// 	// 4. 均不为空 获取指定的
// 	if namespace != "" && name == "" {
// 		result, err = process.DescribeHPA(namespace, name)
// 	}

// 	if err != nil {
// 		fmt.Println("error in process.DescribeHPA ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}

// 	c.JSON(http.StatusOK, result)
// }
