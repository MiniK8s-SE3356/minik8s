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
	"ReplicaSet": getReplicaSet,
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

	// 再获取namespace和name
	namespace, _ := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace, _ = cmd.Flags().GetString("ns")
	}

	// 没有指定namespace，使用default
	if namespace == "" {
		namespace = defaultNamespace
	}

	name := ""
	if len(args) >= 2 {
		name = args[1]
	}

	result, err := getFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	fmt.Println("result is ", result)
}

// get这里有按照namespace获取的功能，这个筛选过程是在前端还是在后端执行?
func getNode(namespace string, name string) (string, error) {
	// 实际上无论namespace和name是什么，getNode都会获取所有的node
	result, err := GetRequest(url.GetNodesURL)
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

	result, err := GetRequestWithParams(url.GetPodsURL, params)
	if err != nil {
		fmt.Println("error in get pod ", err.Error())
		return "", err
	}

	return result, nil
}

func getService(namespace string, name string) (string, error) {
	return "", nil
}

func getReplicaSet(namespace string, name string) (string, error) {
	return "", nil
}

func getNamespace(namespace string, name string) (string, error) {
	return "", nil
}
