package cmdline

import (
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/process"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
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
	"PV":         getPV_PVC,
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

	fmt.Println("result is ", result)
}

func getPV_PVC(namespace string, name string) (string, error) {
	// // 实际上无论namespace和name是什么，getNode都会获取所有的node
	// result, err := httpRequest.GetRequest(RootURL + url.GetNode)
	// if err != nil {
	// 	fmt.Println("error in getNode", err.Error())
	// 	return "", err
	// }

	// return result, nil
	// 请求所有的pv/pvc
	var pv_pvc_list httpobject.HTTPResponse_GetAllPersistVolume
	status, err := httpRequest.GetRequestByObject(RootURL+url.GetAllPersistVolume, nil, &pv_pvc_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return "", err
	}
	writer:=tabwriter.NewWriter(os.Stdout,0,0,2,' ',tabwriter.Debug)
	result:="\n"
	result+="****** PERSISTENT VOLUME ******\n"
	result+="NAME\tTYPE\tCAPACITY\tMOUNTPOD\n"
	for _,pv:=range(pv_pvc_list.Pv){
		mount_str:=""
		for _,mp:=range(pv.Status.MountPod){
			mount_str+=fmt.Sprintf("%s ",mp)
		}
		result+=fmt.Sprintf("%s\t%s\t%s\t%s\n",pv.Metadata.Name,pv.Spec.Type,pv.Spec.Capacity,mount_str)
	}
	fmt.Fprint(writer,result)
	writer.Flush()

	result=""
	result+="****** PERSISTENT VOLUME CLAIM ******\n"
	result+="NAME\tTYPE\tCAPACITY\tBINDPV\n"
	for _,pvc:=range(pv_pvc_list.Pvc){
		bindpv_str:=""
		for _,bpv:=range(pvc.Status.BoundPV){
			bindpv_str+=fmt.Sprintf("%s ",bpv)
		}
		result+=fmt.Sprintf("%s\t%s\t%s\t%s\n",pvc.Metadata.Name,pvc.Spec.Type,pvc.Spec.Capacity,bindpv_str)
	}
	fmt.Fprint(writer,result)
	writer.Flush()
	return  "finish",nil
}

func getNode(namespace string, name string) (string, error) {
	// 实际上无论namespace和name是什么，getNode都会获取所有的node
	result, err := httpRequest.GetRequest(RootURL + url.GetNode)
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

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetPod, params)
	if err != nil {
		fmt.Println("error in get pod ", err.Error())
		return "", err
	}
	formatprint.PrintPods(result)

	return "finish", nil
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
