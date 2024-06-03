package cmdline

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	formatprint "github.com/MiniK8s-SE3356/minik8s/pkg/utils/formatPrint"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/spf13/cobra"
)

var GetFuncTable = map[string]func(namespace string, name string) (string, error){
	"Node":       getNode,
	"Pod":        getPod,
	"Service":    getService,
	"Replicaset": getReplicaSet,
	"HPA":        getHPA,
	"Namespace":  getNamespace,
	"Dns":        getDNS,
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
		namespace = process.DefaultNamespace
	}

	result, err := getFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	//!debug//
	fmt.Println("result is ", result)
	//!debug//
}

func getNode(namespace string, name string) (string, error) {
	// 实际上无论namespace和name是什么，getNode都会获取所有的node
	result, err := httpRequest.GetRequest(RootURL + url.GetNode)
	if err != nil {
		fmt.Println("error in getNode", err.Error())
		return "", err
	}

	formatprint.PrintNodes(result)

	return result, nil
}

func getPod(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetPod, params)
	if err != nil {
		fmt.Println("error in get pod ", err.Error())
		return "", err
	}

	formatprint.PrintPods(result)

	return result, nil
}

func getService(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetService, params)
	if err != nil {
		fmt.Println("error in get service ", err.Error())
		return "", err
	}
	formatprint.PrintService(result)

	return result, nil
}

func getReplicaSet(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetReplicaset, params)
	if err != nil {
		fmt.Println("error in get replicaset ", err.Error())
		return "", err
	}

	formatprint.PrintReplicaset(result)

	return result, nil
}

func getHPA(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetHPA, params)
	if err != nil {
		fmt.Println("error in get HPA ", err.Error())
		return "", err
	}

	formatprint.PrintHPA(result)

	return result, nil
}

func getDNS(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetDNS, params)
	if err != nil {
		fmt.Println("error in get HPA ", err.Error())
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

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetNamespace, params)
	if err != nil {
		fmt.Println("error in get namespace ", err.Error())
		return "", err
	}

	return result, nil
}
