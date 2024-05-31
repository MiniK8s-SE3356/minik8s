package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/gin-gonic/gin"
)

func GetMetricPoint(c *gin.Context) {
	result, err := process.GetMetricPoint()
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, result)
}
