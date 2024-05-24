package cmdline

import (
	"encoding/json"
	"fmt"
	"os"

	minik8s_yaml "github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/yaml"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
)

var applyFuncTable = map[string]func(namespace string, b []byte) error{
	"Pod":        applyPod,
	"Service":    applyService,
	"Replicaset": applyReplicaSet,
	"Namespace":  applyNamespace,
	"HPA":        applyHPA,
}

func ApplyCmdHandler(cmd *cobra.Command, args []string) {
	// 先看一下参数是不是文件路径
	result := checkFilePath(args)
	if !result {
		fmt.Println("not a file")
		cmd.Usage()
		return
	}

	// 读取文件内容，先找到kind
	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("failed to read yaml file")
		return
	}

	var tmp map[string]interface{}
	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		fmt.Println("failed to unmarshal yaml file")
		return
	}

	// kind不支持
	if tmp["kind"] == nil {
		fmt.Println("no kind field found")
		return
	}

	kind := tmp["kind"].(string)
	targetFunc, ok := applyFuncTable[kind]
	if !ok {
		fmt.Println("kind not supported")
		return
	}

	// 再获取namespace
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace, _ = cmd.Flags().GetString("ns")
	}

	// 没有指定namespace，使用default
	if namespace == "" {
		namespace = process.DefaultNamespace
	}

	// 根据kind跳转到相应的处理函数，相当于switch
	err = targetFunc(namespace, data)
	if err != nil {
		fmt.Println(err)
	}
}

func applyPod(namespace string, b []byte) error {
	req := make(map[string]interface{})
	var podDesc minik8s_yaml.PodDesc

	err := yaml.Unmarshal(b, &podDesc)
	if err != nil {
		fmt.Println("failed to unmarshal pod yaml")
		return err
	}
	req["namespace"] = namespace
	req["podDesc"] = podDesc

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("failed to translate into json")
		return err
	}
	// fmt.Println(podDesc.Spec.Containers)
	result, err := httpRequest.PostRequest(url.RootURL+url.AddPod, jsonData)
	if err != nil {
		fmt.Println("error when post request")
		return err
	}

	fmt.Println(result)

	return nil
}

func applyService(namespace string, b []byte) error {
	var ServiceDesc minik8s_yaml.ServiceDesc

	err := yaml.Unmarshal(b, &ServiceDesc)
	if err != nil {
		fmt.Println("failed to unmarshal service yaml")
		return err
	}

	// 这里要区分service下面的子类型
	serviceType, ok := ServiceDesc.Spec.(map[string]interface{})["type"].(string)
	fmt.Println(serviceType)
	if !ok {
		// 默认是ClusterIP
		serviceType = "ClusterIP"
	}

	var jsonData []byte
	var requestMsg struct {
		Namespace string `json:"namespace"`
		Ty        string `json:"ty"`
		Content   string `json:"content"`
	}
	requestMsg.Namespace = namespace
	requestMsg.Ty = serviceType
	requestMsg.Content = string(b)

	// if serviceType == "ClusterIP" {
	// 	var requestMsg struct {
	// 		Namespace string
	// 		Desc      minik8s_yaml.ServiceClusterIPDesc
	// 	}

	// 	var desc minik8s_yaml.ServiceClusterIPDesc
	// 	err := yaml.Unmarshal(b, &desc)
	// 	if err != nil {
	// 		fmt.Println("failed to unmarshal clusterIP service ", err.Error())
	// 		return err
	// 	}

	// 	requestMsg.Desc = desc
	// 	requestMsg.Namespace = namespace

	// 	jsonData, err = json.Marshal(requestMsg)
	// 	if err != nil {
	// 		fmt.Println("failed to marshal clusterIP service", err.Error())
	// 		return err
	// 	}
	// } else if serviceType == "NodePort" {
	// 	var requestMsg struct {
	// 		Namespace string
	// 		Desc      minik8s_yaml.ServiceNodePortDesc
	// 	}

	// 	var desc minik8s_yaml.ServiceNodePortDesc
	// 	err := yaml.Unmarshal(b, &desc)
	// 	if err != nil {
	// 		fmt.Println("failed to unmarshal NodePort service ", err.Error())
	// 		return err
	// 	}

	// 	requestMsg.Desc = desc
	// 	requestMsg.Namespace = namespace

	// 	jsonData, err = json.Marshal(requestMsg)
	// 	if err != nil {
	// 		fmt.Println("failed to marshal NodePort service", err.Error())
	// 		return err
	// 	}
	// }

	fmt.Println(string(jsonData))
	result, err := httpRequest.PostRequest(url.RootURL+url.AddService, jsonData)
	if err != nil {
		fmt.Println("error when post request ", err.Error())
		return err
	}

	fmt.Println(result)

	return nil
}

func applyReplicaSet(namespace string, b []byte) error {
	var request struct {
		Namespace string
		Desc      minik8s_yaml.ReplicaSetDesc
	}
	request.Namespace = namespace

	err := yaml.Unmarshal(b, &request.Desc)
	if err != nil {
		fmt.Println("failed to unmarshal replicaset yaml ", err.Error())
		return err
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return err
	}
	fmt.Println(string(jsonData))

	result, err := httpRequest.PostRequest(url.RootURL+url.AddReplicaset, jsonData)
	if err != nil {
		fmt.Println("error when post request ", err.Error())
		return err
	}

	fmt.Println(result)

	return nil
}

func applyNamespace(namespace string, b []byte) error {
	var namespaceDesc minik8s_yaml.NamespaceDesc

	err := yaml.Unmarshal(b, &namespaceDesc)
	if err != nil {
		fmt.Println("failed to unmarshal namespace yaml")
		return err
	}

	jsonData, err := json.Marshal(namespaceDesc)
	if err != nil {
		fmt.Println("failed to translate into json")
		return err
	}
	result, err := httpRequest.PostRequest(url.RootURL+url.AddNamespace, jsonData)
	if err != nil {
		fmt.Println("error when post request")
		return err
	}

	fmt.Println(result)

	return nil
}

func applyHPA(namespace string, b []byte) error {
	var request struct {
		Namespace string
		Desc      minik8s_yaml.HPADesc
	}
	request.Namespace = namespace

	err := yaml.Unmarshal(b, &request.Desc)
	if err != nil {
		fmt.Println("failed to unmarshal HPA yaml ", err.Error())
		return err
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("failed to translate into json ", err.Error())
		return err
	}
	fmt.Println(string(jsonData))

	result, err := httpRequest.PostRequest(url.RootURL+url.AddHPA, jsonData)
	if err != nil {
		fmt.Println("error when post request ", err.Error())
		return err
	}

	fmt.Println(result)

	return nil

}

func checkFilePath(args []string) bool {
	// 检查参数给出的文件路径是否存在

	if len(args) == 0 {
		return false
	}

	result, err := os.Stat(args[0])
	if err != nil {
		return false
	}

	if result.IsDir() {
		return false
	}

	return true
}
