package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"

	"github.com/gin-gonic/gin"
)

func AddPod(c *gin.Context) {
	var desc struct {
		PodDesc   yaml.PodDesc `json:"podDesc"`
		Namespace string       `json:"namespace"`
	}
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddPod(process.DefaultNamespace, &desc.PodDesc)
	if err != nil {
		fmt.Println("error in process.AddPod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func RemovePod(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// namespace, ok := param["namespace"]
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'namespace'"})
	// 	return
	// }

	name, ok := param["name"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'name'"})
		return
	}

	result, err := process.RemovePod(process.DefaultNamespace, name)
	if err != nil {
		fmt.Println("error in process.RemovePod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// apply是不是已经够了，不需要modify？
// func ModifyPod(c *gin.Context) {

// }

func UpdatePod(c *gin.Context) {
	var req struct {
		Pod       pod.Pod `json:"pod"`
		Namespace string  `json:"namespace"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.UpdatePod(process.DefaultNamespace, req.Pod)

	if err != nil {
		fmt.Println("error in process updatePod")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetPod(c *gin.Context) {
	// var param map[string]string
	// if err := c.ShouldBindJSON(&param); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	namespace := c.Query("namespace")
	name := c.Query("name")

	var result map[string]interface{}
	var err error
	// 四种情况
	// 1. namespace name均为空 获取全部pod
	if namespace == "" && name == "" {
		result, err = process.GetAllPods()
	}
	// 2. namespace为空，name不为空 获取Default下的pod
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.GetPod(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部pod
	if namespace != "" && name == "" {
		result, err = process.GetPods(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name != "" {
		result, err = process.GetPod(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.GetPod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func GetAllPod(c *gin.Context) {
	result, err := process.GetAllPods()

	if err != nil {
		fmt.Println("error in process.DescribePod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func DescribePod(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// namespace := c.Query("namespace")
	namespace := process.DefaultNamespace

	name := c.Query("name")

	var result string
	var err error
	// 四种情况
	// 1. namespace name均为空 获取全部pod
	if namespace == "" && name == "" {
		result, err = process.DescribeAllPods()
	}
	// 2. namespace为空，name不为空 获取Default下的pod
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err = process.DescribePod(namespace, name)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部pod
	if namespace != "" && name == "" {
		result, err = process.DescribePods(namespace)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name == "" {
		result, err = process.DescribePod(namespace, name)
	}

	if err != nil {
		fmt.Println("error in process.DescribePod ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// func GetPods(c *gin.Context) {
// 	var param map[string]string
// 	if err := c.ShouldBindJSON(&param); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	namespace, ok := param["namespace"]
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'namespace'"})
// 		return
// 	}

// 	result, err := process.GetPods(namespace)
// 	if err != nil {
// 		fmt.Println("error in process.GetPods ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}

// 	c.JSON(http.StatusOK, result)
// }
