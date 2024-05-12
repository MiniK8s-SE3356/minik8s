package cmdline

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/spf13/cobra"
)

var deleteFuncTable = map[string]func(namespace string, name string) (string, error){
	"Pod":        deletePod,
	"Service":    deleteService,
	"ReplicaSet": deleteReplicaSet,
	"Namespace":  deleteNamespace,
}

func DeleteCmdHandler(cmd *cobra.Command, args []string) {
	// 整体的逻辑先按get移植过来
	if len(args) == 0 {
		cmd.Usage()
		return
	}

	// 先获取kind
	kind := args[0]
	deleteFunc, ok := deleteFuncTable[kind]
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

	if kind == "Namespace" && name == "" {
		fmt.Println("name of the namespace not found")
		return
	}

	result, err := deleteFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	fmt.Println("result is ", result)
}

func deletePod(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RemovePodURL, params)
	if err != nil {
		fmt.Println("error in delete pod ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteService(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RemoveServiceURL, params)
	if err != nil {
		fmt.Println("error in delete service ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteReplicaSet(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := GetRequestWithParams(url.RemoveReplicasetURL, params)
	if err != nil {
		fmt.Println("error in delete replicaset ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteNamespace(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": name,
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}

	result, err := PostRequest(url.RemoveNamespaceURL, jsonData)
	if err != nil {
		fmt.Println("error in delete namespace ", err.Error())
		return "", err
	}

	return result, nil
}
