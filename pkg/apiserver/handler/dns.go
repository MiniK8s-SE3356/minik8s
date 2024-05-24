package handler

import (
	"fmt"
	"net/http"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"

	"github.com/gin-gonic/gin"
)

// gin server的 /api/v1/addDNS对应的方法
func AddDNS(c *gin.Context) {
	var requestMsg struct {
		Namespace string
		Desc      yaml.DNSDesc
	}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.AddDNS(requestMsg.Namespace, &requestMsg.Desc)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RemoveDNS(c *gin.Context) {
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

	result, err := process.RemoveDNS(namespace, name)
	if err != nil {
		fmt.Println("error in process.RemoveDNS ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func UpdateDNS(c *gin.Context) {
	var requestMsg struct {
		Namespace string
		dns       dns.DNS
	}
	if err := c.ShouldBindJSON(&requestMsg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := process.UpdateDNS(requestMsg.Namespace, requestMsg.dns)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetDNS(c *gin.Context) {
	// namespace := c.Query("namespace")
	namespace := process.DefaultNamespace
	name := c.Query("name")

	// 四种情况
	// 1. namespace name均为空 获取全部dns
	if namespace == "" && name == "" {
		result, err := process.GetAllDNSs()

		if err != nil {
			fmt.Println("error in process.GetDNS ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 2. namespace为空，name不为空 获取Default下的dns
	if namespace == "" && name != "" {
		namespace = "Default"
		result, err := process.GetDNS(namespace, name)

		if err != nil {
			fmt.Println("error in process.GetDNS ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 3. namespace不为空，name为空 获取namespace下的全部dns
	if namespace != "" && name == "" {
		result, err := process.GetDNSs(namespace)

		if err != nil {
			fmt.Println("error in process.GetDNS ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}
	// 4. 均不为空 获取指定的
	if namespace != "" && name != "" {
		result, err := process.GetDNS(namespace, name)

		if err != nil {
			fmt.Println("error in process.GetDNS ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, result)
	}

}

func GetAllDNS(c *gin.Context) {
	result, err := process.GetAllDNSs()
	if err != nil {
		fmt.Println("error in process.GetAllDNSs ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

// func DescribeDNS(c *gin.Context) {
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
// 	// 1. namespace name均为空 获取全部dns
// 	if namespace == "" && name == "" {
// 		result, err = process.DescribeAllDNSs()
// 	}
// 	// 2. namespace为空，name不为空 获取Default下的dns
// 	if namespace == "" && name != "" {
// 		namespace = "Default"
// 		result, err = process.DescribeDNS(namespace, name)
// 	}
// 	// 3. namespace不为空，name为空 获取namespace下的全部dns
// 	if namespace != "" && name == "" {
// 		result, err = process.DescribeDNSs(namespace)
// 	}
// 	// 4. 均不为空 获取指定的
// 	if namespace != "" && name == "" {
// 		result, err = process.DescribeDNS(namespace, name)
// 	}

// 	if err != nil {
// 		fmt.Println("error in process.DescribeDNS ", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}

// 	c.JSON(http.StatusOK, result)
// }
