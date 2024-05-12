package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/types/node"

	"github.com/gin-gonic/gin"
)

func AddNode(c *gin.Context) {
	var desc node.Node
	if err := c.ShouldBindJSON(&desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddNode(&desc)
	if err != nil {
		fmt.Println("error in process.AddNode ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func RemoveNode(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name, ok := param["name"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no field 'name'"})
		return
	}

	result, err := process.RemoveNode(name)
	if err != nil {
		fmt.Println("error in process.RemoveNode ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// apply是不是已经够了，不需要modify？
// func ModifyNode(c *gin.Context) {

// }

func GetNode(c *gin.Context) {
	var param map[string]string
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := c.Query("name")

	if name == "" {
		result, err := process.GetNodes()

		if err != nil {
			fmt.Println("error in process.GetNodes")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	} else {
		// node是没有namespace的
		result, err := process.GetNode(name)

		if err != nil {
			fmt.Println("error in process.GetNode ")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
}

// func DescribeNode(c *gin.Context) {
// 	var param map[string]string
// 	if err := c.ShouldBindJSON(&param); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	name := c.Query("name")

// 	result, err := process.DescribeNode(name)

// 	if err != nil {
// 		fmt.Println("error in process.DescribeNode ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}

// 	c.JSON(http.StatusOK, result)
// }
