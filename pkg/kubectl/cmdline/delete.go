package cmdline

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/spf13/cobra"
)

var deleteFuncTable = map[string]func(namespace string, name string) (string, error){
	"Pod":        deletePod,
	"Service":    deleteService,
	"Replicaset": deleteReplicaSet,
	"HPA":        deleteHPA,
	"Namespace":  deleteNamespace,
	"Dns":        deleteDNS,
	"Function":   deleteFunction,
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
		namespace = process.DefaultNamespace
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
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemovePod, jsonData)
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
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemoveService, jsonData)
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

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemoveReplicaset, jsonData)
	if err != nil {
		fmt.Println("error in delete replicaset ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteHPA(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemoveHPA, jsonData)
	if err != nil {
		fmt.Println("error in delete HPA ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteDNS(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemoveDNS, jsonData)
	if err != nil {
		fmt.Println("error in delete HPA ", err.Error())
		return "", err
	}

	return result, nil
}

func deleteFunction(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Println("failed to translate into json")
		return "", err
	}
	result, err := httpRequest.PostRequest(RootURL+url.RemoveFunction, jsonData)
	if err != nil {
		fmt.Println("error in delete Function ", err.Error())
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

	result, err := httpRequest.PostRequest(RootURL+url.RemoveNamespace, jsonData)
	if err != nil {
		fmt.Println("error in delete namespace ", err.Error())
		return "", err
	}

	return result, nil
}
