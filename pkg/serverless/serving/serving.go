package serving

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/pod"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiserver/url"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/config"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
)

func ScaleFunctionPod() {
	frequency := config.GetFuncionPodRequestFrequency()
	fmt.Println(frequency)
	for funcName, f := range frequency {
		fmt.Println("function name: ", funcName, " frequency: ", f)
		if f > 3 {
			scaleUp(funcName)
		} else if f < 0.01 {
			scaleDown(funcName)
		}
	}
}

func scaleUp(funcName string) {
	var req struct {
		FuncName string `json:"FuncName"`
	}
	req.FuncName = funcName
	jsonData, _ := json.Marshal(req)
	result, err := httpRequest.PostRequest(config.HTTPURL_root+url.AddServerlessFuncPod, jsonData)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result, err)

}

func scaleDown(funcName string) {
	params := make(map[string]string)
	params["funcName"] = funcName
	result, err := httpRequest.GetRequestWithParams(config.HTTPURL_root+url.GetServerlessFuncPod, params)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pods []pod.Pod
	json.Unmarshal([]byte(result), &pods)

	if len(pods) != 0 {
		var req struct {
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		}
		req.Namespace = pods[0].Metadata.Namespace
		req.Name = pods[0].Metadata.Name

		jsonData, _ := json.Marshal(req)

		result, err := httpRequest.PostRequest(config.HTTPURL_root+url.RemovePod, jsonData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(result)
	}
}
