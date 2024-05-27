package hostsController

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/dns"
	"github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/config"
	kptype "github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/types"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/txn2/txeh"
)

type HostsController struct {
	// nginx所在的节点（一般为控制平面）的IP
	NginxIp string
	// Hosts *txeh.Hosts
}

func NewHostsController() *HostsController {
	return &HostsController{}
}

func (hc *HostsController) Init() {
	// TODO: 获取真实且动态的Nginx所在节点Ip，目前只是写死的
	// hosts, err := txeh.NewHostsDefault()
	// if err != nil {
	// 	fmt.Printf("Can't create hosts by txeh\n")
	// 	panic(err)
	// }
	hc.NginxIp = config.NGINX_IP
	// hc.Hosts=hosts
}

// func (hc *HostsController) Run() {
// 	poller.PollerStaticPeriod(10*time.Second, hc.routine, true)
// }

func (hc *HostsController) SyncEtcHosts(dns_list *httpobject.HTTPResponse_GetAllDns) (kptype.KpDnsStatus, error) {
	result := kptype.KpDnsStatus{}
	// 创建新的dns hosts
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		fmt.Printf("Can't create NewHostsDefault!\n")
		return kptype.KpDnsStatus{}, err
	}

	for _, dns_item := range *dns_list {
		if dns_item.Status.Phase != dns.DNS_READY {
			continue
		}
		// /etc/hosts加入的条目为：所有dns的host(不含port和子路径)到nginxIP的本地域名解析规则
		hosts.AddHost(hc.NginxIp, dns_item.Spec.Host)
		// 待返回结果更新
		result[dns_item.Metadata.Name] = kptype.KpDns{
			Version: dns_item.Status.Version,
		}
	}

	// 更新的结构存入/etc/hosts中
	err = hosts.Save()
	if err != nil {
		fmt.Printf("Can't Update /etc/hosts!\n")
		return kptype.KpDnsStatus{}, err
	}

	return result, nil
}
