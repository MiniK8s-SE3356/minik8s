package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	version     = "v1"
	getPodsURL  = "/api/" + version + "/getPods"
	getNodesURL = "/api/" + version + "/getNodes"
)

func example(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func bind(r *gin.Engine) {
	r.GET(getPodsURL, example)
	r.GET(getNodesURL, example)
}

func main() {
	r := gin.Default()
	bind(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
