package cmdline

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/spf13/cobra"
)

var GetFuncTable = map[string]func(namespace string, name string) (string, error){
	"Node":       getNode,
	"Pod":        getPod,
	"Service":    getService,
	"Replicaset": getReplicaSet,
	"Namespace":  getNamespace,
}

func GetCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	// 先获取kind
	kind := args[0]
	getFunc, ok := GetFuncTable[kind]
	if !ok {
		fmt.Println("kind not supported")
		return
	}

	name := ""
	if len(args) >= 2 {
		name = args[1]
	}

	// 再获取namespace
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace, _ = cmd.Flags().GetString("ns")
	}

	// 没有指定namespace，使用default
	if name != "" && namespace == "" {
		namespace = defaultNamespace
	}

	result, err := getFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	fmt.Println("result is ", result)
}

func getNode(namespace string, name string) (string, error) {
	// 实际上无论namespace和name是什么，getNode都会获取所有的node
	result, err := GetRequest(url.RootURL + url.GetNode)
	if err != nil {
		fmt.Println("error in getNode", err.Error())
		return "", err
	}

	return result, nil
}

func getPod(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.GetPod, params)
	if err != nil {
		fmt.Println("error in get pod ", err.Error())
		return "", err
	}

	return result, nil
}

func getService(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.GetService, params)
	if err != nil {
		fmt.Println("error in get service ", err.Error())
		return "", err
	}

	return result, nil
}

func getReplicaSet(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.GetReplicaset, params)
	if err != nil {
		fmt.Println("error in get replicaset ", err.Error())
		return "", err
	}

	return result, nil
}

func getNamespace(namespace string, name string) (string, error) {
	// 这里用name传
	params := map[string]string{
		// "namespace": namespace,
		"name": name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.GetNamespace, params)
	if err != nil {
		fmt.Println("error in get namespace ", err.Error())
		return "", err
	}

	return result, nil
}
