package dnsController

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	"github.com/MiniK8s-SE3356/minik8s/pkg/controller/dnsController/nginxManager"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/httpRequest"
	"github.com/MiniK8s-SE3356/minik8s/pkg/utils/poller"
)

type DnsController struct {
	nginxManager *nginxManager.NginxManager
}

func NewDnsController() *(DnsController) {
	fmt.Printf("New DnsController\n")
	nginx_manager := nginxManager.NewDnsManager()
	return &DnsController{
		nginxManager: nginx_manager,
	}
}

func (dc *DnsController) Init() {
	fmt.Printf("Init DnsController\n")
	dc.nginxManager.Init()
}

func (dc *DnsController) Run() {
	fmt.Printf("Run DnsController\n")
	poller.PollerStaticPeriod(5*time.Second, dc.routine, true)

}

// func is_exits(spec *dns.DnsPathInfo,status_list *map[string]dns.DnsPathStatus)bool{
// 	_,exist:=(*status_list)[spec.SubPath]
// 	return exist
// }

func (dc *DnsController) routine() {
	fmt.Printf("DnsController routine ...\n")

	var service_list httpobject.HTTPResponse_GetAllServices = httpobject.HTTPResponse_GetAllServices{}
	var dns_list httpobject.HTTPResponse_GetAllDns = httpobject.HTTPResponse_GetAllDns{}
	var dns_update_list httpobject.HTTPRequest_UpdateDns = httpobject.HTTPRequest_UpdateDns{}
	var clusterip_name_map map[string]*service.ClusterIP = make(map[string]*service.ClusterIP)
	var wg sync.WaitGroup

	// 获得所有service
	status, err := httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllService", nil, &service_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		goto return_directly
	}

	// 在正常请求下所有的数据后，一方面需要使用dns更新nginx,另一方面需要dns和service共同更新dns
	// 我们将这两件事情并行去做，并阻塞式等待两者全部完成再退出函数

	// 获得所有dns
	status, err = httpRequest.GetRequestByObject("http://192.168.1.6:8080/api/v1/GetAllDNS", nil, &dns_list)
	if status != http.StatusOK || err != nil {
		fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
		goto return_directly
	}

	fmt.Println(service_list.ClusterIP)
	fmt.Println(dns_list)

	// 根据service更新dns
	// 作为工作队列中的协程，此线程会在退出时通知主线程，主线程也会在最后阻塞等待所有协程return
	wg.Add(1)
	go dc.syncNginx(&(dns_list), &wg)

	// 针对clusterip建立name->point的map表
	for _, clusterip_item := range service_list.ClusterIP {
		//没有READY的clusterip不加入
		if clusterip_item.Status.Phase != service.CLUSTERIP_READY {
			continue
		}
		clusterip_name_map[clusterip_item.Metadata.Name] = &clusterip_item
	}

	// fmt.Println(clusterip_name_map)

	// 根据上述的clusterip数据结构更新dns,并将之前未READY,在本轮更新后READY的dns加入更新列表
	for _, dns_item := range dns_list {
		need_update := false
		if(dns_item.Status.PathsStatus==nil){
			need_update=true
			dns_item.Status.PathsStatus=make(map[string]dns.DnsPathStatus)
		}
		// 目前的READY策略是，DNS只要有了下属的CLUSTERIP,就可以READY,后续状态会逐渐更新，version也会随之提升
		// 如果后续发生dns对象status更改，则会设置need_update为true
		
		// 遍历所有的spec状态
		for _, port_spec_item := range dns_item.Spec.Paths {
			// 1.1 spec的clusterip存在
			//		2.1 spec对应的项在status中不存在	-> 需要添加status，need update
			//		2.2 spec对应的项在status中存在
			//			3.1 status和spec不符		-> 需要更新status，need update
			//			3.2 status和spec相符		-> 无需更新status
			// 1.2spec的clusterip不存在或未READY
			// 		2.3 spec对应的项在status中不存在	-> 无需更新status
			//		2.4 spec对应的项在status中存在		-> 需要删除status, need update

			clusterip_ptr, clusterip_exist := clusterip_name_map[port_spec_item.SvcName]
			status_value, status_exist := dns_item.Status.PathsStatus[port_spec_item.SvcName]
			if clusterip_exist {
				if status_exist {
					if status_value.SvcIP != clusterip_ptr.Metadata.Ip {
						// 3.1
						need_update = true
						dns_item.Status.PathsStatus[port_spec_item.SvcName] = dns.DnsPathStatus{
							SubPath: port_spec_item.SubPath,
							SvcIP:   clusterip_ptr.Metadata.Ip,
							SvcPort: port_spec_item.SvcPort,
						}
					} else {
						// 3.2
					}
				} else {
					// 2.1
					need_update = true
					dns_item.Status.PathsStatus[port_spec_item.SvcName] = dns.DnsPathStatus{
						SubPath: port_spec_item.SubPath,
						SvcIP:   clusterip_ptr.Metadata.Ip,
						SvcPort: port_spec_item.SvcPort,
					}
				}
			} else {
				if status_exist {
					// 2.4
					need_update = true
					delete(dns_item.Status.PathsStatus, port_spec_item.SvcName)
				} else {
					// 2.3
				}
			}

		}

		if need_update {
			if len(dns_item.Status.PathsStatus) > 0 {
				dns_item.Status.Phase = dns.DNS_READY
			}
			dns_item.Status.Version += 1
			dns_update_list[dns_item.Metadata.Name] = dns_item
		}
	}

	fmt.Println(dns_update_list)

	if len(dns_update_list) > 0 {
		// 将需要更新的dns更新回去
		status, err = httpRequest.PostRequestByObject("http://192.168.1.6:8080/api/v1/UpdateDNS", dns_update_list, nil)
		if status != http.StatusOK || err != nil {
			fmt.Printf("EndpointsController routine error get, status %d, return\n", status)
			goto return_with_sync_wokequeue
		}
	}

	// 阻塞等待所有协程完成执行
return_with_sync_wokequeue:
	wg.Wait()
return_directly:
	return
}

func (dc *DnsController) syncNginx(dns_list *httpobject.HTTPResponse_GetAllDns, wg *sync.WaitGroup) {
	defer wg.Done()

	dc.nginxManager.SyncNginx(dns_list)
}
