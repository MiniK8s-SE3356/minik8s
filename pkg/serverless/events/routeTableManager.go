package events

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/selector"
	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/config"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/selectorUtils"
)

const (
	ROUTETABLE_NONEPOD = "NONE"
)

type routeInfo struct {
	next      int
	podIPList []string
}

// RouteTableManagers是**线程安全**的
type RouteTableManager struct {
	/* 路由表，function name -> pod IP list */
	routeTable      map[string]routeInfo
	routeTableMutex sync.Mutex
}

func NewRouteTableManager() *RouteTableManager {
	fmt.Printf("New RouteTableManager...\n")
	return &RouteTableManager{}
}

func (rm *RouteTableManager) Init() {
	fmt.Printf("Init RouteTableManager...\n")
	rm.routeTable = make(map[string]routeInfo)
	rm.routeTableMutex = sync.Mutex{}
}

func (rm *RouteTableManager) SyncRoutine() {
	fmt.Printf("Sync RouteTable...\n")
	// 请求所有的pod
	var pod_list httpobject.HTTPResponse_GetAllPod
	status, err := httpRequest.GetRequestByObject(config.HTTPURL_GetAllPod, nil, &pod_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get,  status %d, return\n", status)
		return
	}

	// 请求所有的serverless function name
	var serverless_func_list httpobject.HTTPResponse_GetAllServerlessFunction
	status, err = httpRequest.GetRequestByObject(config.HTTPURL_GetAllServerlessFunction, nil, &serverless_func_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("routine error get, status %d, return\n", status)
		return
	}

	// 创建new_route_table
	new_route_table := make(map[string]routeInfo)
	// 创建Selector
	func_pod_selector := selector.Selector{
		MatchLabels: make(map[string]string),
	}

	// 为new_route_table填写内容
	for _, serverless_func_name := range serverless_func_list {
		// TODO：具体的标签筛选规则待商议
		// 修改筛选器
		func_pod_selector.MatchLabels["serverlessFuncName"] = serverless_func_name
		// 复用之前的筛选工具，选出pod name list
		new_f2p := selectorUtils.SelectPodNameList(&func_pod_selector, &pod_list)
		// 根据pod name list,转化为pod ip list（允许为空）
		new_podIP_list := []string{}
		for _, pod_name := range new_f2p {
			// FIXME: 如果podIP是null,会不会引发问题？
			// FIXME: podIP是否可用，是否还要看pod的其他状态？
			pip := pod_list[pod_name].Status.PodIP
			if pip != "" {
				new_podIP_list = append(new_podIP_list, pip)
			}
		}
		// 为new_route_table填写一条条目
		new_route_table[serverless_func_name] = routeInfo{
			next:      0,
			podIPList: new_podIP_list,
		}
	}

	// 如果new_route_table正常，则用于更新
	rm.routeTableMutex.Lock()
	rm.routeTable = new_route_table
	// println(rm.routeTable)
	fmt.Println("rm.routeTable: ", rm.routeTable)
	rm.routeTableMutex.Unlock()
}

func (rm *RouteTableManager) FunctionName2PodIP(funcName string) string {
	result := ROUTETABLE_NONEPOD
	rm.routeTableMutex.Lock()
	if va, ex := rm.routeTable[funcName]; ex {
		if len(va.podIPList) <= 0 {
			goto return_with_unlock
		}
		va.next = (va.next + 1) % (len(va.podIPList))
		result = va.podIPList[va.next]
	}
return_with_unlock:
	rm.routeTableMutex.Unlock()
	return result
}

func (rm *RouteTableManager) GetFunctionPodNum(funcName string) int {
	result := 0
	rm.routeTableMutex.Lock()
	if va, ex := rm.routeTable[funcName]; ex {
		result = len(va.podIPList)
	}
	rm.routeTableMutex.Unlock()
	return result
}
