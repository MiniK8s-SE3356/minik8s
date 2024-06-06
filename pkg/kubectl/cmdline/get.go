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
	"Function":   getFunction,
	"GPUJob":     getGPUJob,
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

	_, err := getFunc(namespace, name)
	if err != nil {
		fmt.Println("error in GetCmdHandler ", err.Error())
		return
	}

	//!debug//
	// fmt.Println("result is ", result)
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

	return "", nil
}

func getService(namespace string, name string) (string, error) {
	// params := map[string]string{
	// 	"namespace": namespace,
	// 	"name":      name,
	// }

	// result, err := httpRequest.GetRequestWithParams(RootURL+url.GetService, params)
	// if err != nil {
	// 	fmt.Println("error in get service ", err.Error())
	// 	return "", err
	// }
	// return result, nil
	// 获得所有service
	var service_list httpobject.HTTPResponse_GetAllServices = httpobject.HTTPResponse_GetAllServices{}
	status, err := httpRequest.GetRequestByObject(RootURL+url.GetAllService, nil, &service_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		return "", err
	}
	// 获得所有endpoint
	var endpoint_list httpobject.HTTPResponse_GetAllEndpoint
	status, err = httpRequest.GetRequestByObject(RootURL+url.GetAllEndpoint, nil, &endpoint_list)

	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return "", err
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

	// 先组装ClusterIP格式化展示字符串
	result := "\n"
	result += "****** CLUSTERIP ******\n"
	result += "NAME\tSELECTOR\tIP\tPORTS\tENDPOINTS\n"
	for _, clu_item := range service_list.ClusterIP {
		selector_str := ""
		ports_str := ""
		endponints_str := ""
		// 组装selector字符串
		for k, v := range clu_item.Spec.Selector.MatchLabels {
			selector_str += fmt.Sprintf("%s:%s ", k, v)
		}

		// 组装ports字符串
		for _, v := range clu_item.Spec.Ports {
			ports_str += fmt.Sprintf("%d:%d ", int(v.Port), int(v.TargetPort))
		}

		// 组装endpoints字符串
		for _, e_list := range clu_item.Status.ServicesStatus {
			for _, e_id := range e_list {
				if v, exist := endpoint_list[e_id]; exist {
					endponints_str += fmt.Sprintf("%s:%d ", v.PodIP, int(v.PodPort))
				}
			}
		}

		final_str := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", clu_item.Metadata.Name, selector_str, clu_item.Metadata.Ip, ports_str, endponints_str)
		result += final_str
	}
	fmt.Fprint(writer, result)
	writer.Flush()

	// 后组装NodePort格式化展示字符串
	result = "\n"
	result += "****** NODEPORT ******\n"
	result += "NAME\tSELECTOR\tBINDCLUSTERIP\tPORTS\n"
	for _, np_item := range service_list.NodePort {
		selector_str := ""
		// 组装selector字符串
		for k, v := range np_item.Spec.Selector.MatchLabels {
			selector_str += fmt.Sprintf("%s:%s ", k, v)
		}
		ports_str := ""
		// 组装ports字符串
		for _, v := range np_item.Spec.Ports {
			ports_str += fmt.Sprintf("%d->%d:%d ", int(v.NodePort), int(v.Port), int(v.TargetPort))
		}

		final_str := fmt.Sprintf("%s\t%s\t%s\t%s\n", np_item.Metadata.Name, selector_str, np_item.Status.ClusterIPID, ports_str)
		result += final_str
	}

	fmt.Fprint(writer, result)
	writer.Flush()

	return "finish", nil
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
	// params := map[string]string{
	// 	"namespace": namespace,
	// 	"name":      name,
	// }

	// result, err := httpRequest.GetRequestWithParams(RootURL+url.GetDNS, params)
	// if err != nil {
	// 	fmt.Println("error in get HPA ", err.Error())
	// 	return "", err
	// }

	// return result, nil

	var dns_list httpobject.HTTPResponse_GetAllDns = httpobject.HTTPResponse_GetAllDns{}
	status, err := httpRequest.GetRequestByObject(RootURL+url.GetAllDNS, nil, &dns_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		return "", err
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	result := "\n"
	result += "****** DNS ******\n"
	result += "NAME\tKIND\tHOST\tPATHS\n"
	for _, dns_item := range dns_list {
		path_str := ""
		// 组装path_str
		for _, v := range dns_item.Status.PathsStatus {
			path_str += fmt.Sprintf("%s->%s:%d ", v.SubPath, v.SvcIP, int(v.SvcPort))
		}
		result += fmt.Sprintf("%s\t%s\t%s\t%s\n", dns_item.Metadata.Name, dns_item.Kind, dns_item.Spec.Host, path_str)
	}

	fmt.Fprint(writer, result)
	writer.Flush()
	return "finish", nil
}

func getGPUJob(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(GPUCtlRootURL+url.GetGPUJob, params)
	if err != nil {
		fmt.Println("error in get GPUJob ", err.Error())
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

func getFunction(namespace string, name string) (string, error) {
	params := map[string]string{
		"namespace": namespace,
		"name":      name,
	}

	result, err := httpRequest.GetRequestWithParams(RootURL+url.GetFunction, params)
	if err != nil {
		fmt.Println("error in get Function ", err.Error())
		return "", err
	}

	formatprint.PrintFunction(result)
	return result, nil
}
