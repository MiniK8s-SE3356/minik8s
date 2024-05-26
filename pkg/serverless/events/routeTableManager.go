package events

import (
	"fmt"
	"net/http"
	"sync"

	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
)

const (
	ROUTETABLE_NODEPOD = "NONE"
)

var routeTableMutex sync.Mutex

type routeInfo struct {
	next      int
	podIPList []string
}

type RouteTableManager struct {
	/* 路由表，function name -> pod IP list */
	routeTable map[string]routeInfo
}

func NewRouteTableManager() *RouteTableManager {
	fmt.Printf("New RouteTableManager...\n")
	return &RouteTableManager{
		routeTable: make(map[string]routeInfo),
	}
}

func (rm *RouteTableManager) Init() {
	fmt.Printf("Init RouteTableManager...\n")
	routeTableMutex = sync.Mutex{}
}

func (rm *RouteTableManager) SyncRoutine() {
	fmt.Printf("Sync RouteTable...\n")
	// 请求所有的pod
	var pod_list httpobject.HTTPResponse_GetAllPod
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllPod", nil, &pod_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("Sync RouteTable Routine error Get, status %d, return\n", status)
		return
	}
	// 创建new_route_table
	new_route_table:=make(map[string]routeInfo)

	// 如果new_route_table正常，则用于更新
	routeTableMutex.Lock()
	rm.routeTable=new_route_table
	routeTableMutex.Unlock()
}

func (rm *RouteTableManager) FunctionName2PodIP(funName string) string {
	result := ROUTETABLE_NODEPOD
	if va, ex := rm.routeTable[funName]; ex {
		routeTableMutex.Lock()
		if len(va.podIPList) <= 0 {
			goto return_with_unlock
		}
		va.next = (va.next + 1) % (len(va.podIPList))
		result = va.podIPList[va.next]
	return_with_unlock:
		routeTableMutex.Unlock()
		return result
	}
	return result
}
