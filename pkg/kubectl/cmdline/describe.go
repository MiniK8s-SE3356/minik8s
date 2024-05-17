package cmdline

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/spf13/cobra"
)

var describeFuncTable = map[string]func(namespace string, name string) (string, error){
	"Pod":        describePod,
	"Service":    describeService,
	"ReplicaSet": describeReplicaSet,
	"Namespace":  describeNamespace,
}

func DescribeCmdHandler(cmd *cobra.Command, args []string) {
	// 整体的逻辑先按get移植过来
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	// 先获取kind
	kind := args[0]
	describeFunc, ok := describeFuncTable[kind]
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

	result, err := describeFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	fmt.Println("result is ", result)
}

func describePod(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.DescribePod, params)
	if err != nil {
		fmt.Println("error in describe pod ", err.Error())
		return "", err
	}

	return result, nil
}

func describeService(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.DescribeService, params)
	if err != nil {
		fmt.Println("error in describe service ", err.Error())
		return "", err
	}

	return result, nil
}

func describeReplicaSet(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.DescribeReplicaset, params)
	if err != nil {
		fmt.Println("error in describe replicaset ", err.Error())
		return "", err
	}

	return result, nil
}

func describeNamespace(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		// "name":      name,
	}

	result, err := GetRequestWithParams(url.RootURL+url.DescribeNamespace, params)
	if err != nil {
		fmt.Println("error in describe namespace ", err.Error())
		return "", err
	}

	return result, nil
}
