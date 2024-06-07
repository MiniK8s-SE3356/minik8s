package iptablesController

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MiniK8s-SE3356/minik8s/pkg/apiObject/service"
	kptype "github.com/MiniK8s-SE3356/minik8s/pkg/kubeProxy/types"
	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
	"github.com/coreos/go-iptables/iptables"
)

const (
	KUBE_SERVICES   = "KUBE-SERVICE"
	KUBE_NODEPORTS  = "KUBE-NODEPORTS"
	KUBE_MARK_MASQ  = "KUBE-MARK-MASQ"
	KUBE_POSTROUING = "KUBE-POSTROUING"
	KUBE_SVC        = "KUBE-SVC-"
	KUBE_SEP        = "KUBE-SEP-"
	POSTROUING      = "POSTROUTING"
	PREROUTING      = "PREROUTING"
	OUTPUT          = "OUTPUT"
	INPUT           = "INPUT"
	DNAT            = "DNAT"

	not_found     = "not found"
	already_exist = "already exist"
	add_success   = "add success"
)

type metadata struct {
	ip       string
	protocol string /* 只允许tcp和udp两种字符串 */
	port     uint16
}

type IptablesController struct {
	// 控制iptables的实际对象
	ipt *iptables.IPTables

	// 记录service+port下辖pods，service ID+port --> endpoint ID数组
	// ID是该类对象的集群唯一标识符，这些符号会参与构建KUBE-SVC-XXX和KUBE-SEP-XXX
	sid2eid map[string]([]string)

	// 记录nodePort与service+port对应关系，本机端口号 --> service ID+port
	nodeport2sid map[uint16](string)

	// service ID+port --> service metadata的map
	// service metadata记录了service的虚拟IP，协议（tcp/udp）和虚拟port
	sid2smeta map[string]metadata

	// endpoint ID --> endpoint metadata的map
	// endpoint metadata记录了endpoint被pod IP,协议（tcp/udp）和port
	eid2emeta map[string]metadata
}

// NewIptablesController 创建一个新的IptablesController对象，但是没有任何初始化工作
//
//	@return *IptablesController
func NewIptablesController() *IptablesController {
	fmt.Printf("Create IptablesController ...\n")
	// 创建IptablesController
	iptables_controller := &IptablesController{}
	return iptables_controller
}

// Init 初始化IptablesController对象，并重置相关iptables规则为初始状态（所有的KUBE-SVC-XXX和KUBE-SEP-XXX都不存在，其他KUBE链均配好初始状态且插入到对应的原装链位置）
//
//	@receiver ic
func (ic *IptablesController) Init() {
	fmt.Printf("Init IptablesController ...\n")
	// 创建iptables实例
	ic.ipt, _ = iptables.New()
	// 为所有的map做初始化
	ic.sid2eid = make(map[string][]string)
	ic.nodeport2sid = make(map[uint16]string)
	ic.sid2smeta = make(map[string]metadata)
	ic.eid2emeta = make(map[string]metadata)

	// Naive trail for iptables, including handling error
	// err :=ic.ipt.Append("nat", "PREROUTING", "-j", "ACCEPT")
	// if(err!=nil){
	// 	fmt.Printf("Error msg: %s",err.Error())
	// }

	// 以下是对nat table的修改，目前只需要修改此表
	// 设置NAT表基本链的默认策略全部为ACCEPT
	ic.ipt.ChangePolicy("nat", PREROUTING, "ACCEPT")
	ic.ipt.ChangePolicy("nat", INPUT, "ACCEPT")
	ic.ipt.ChangePolicy("nat", OUTPUT, "ACCEPT")
	ic.ipt.ChangePolicy("nat", POSTROUING, "ACCEPT")

	// KUBE_SERVICES和KUBE_NODEPORTS：若有则清空，若无则新建
	var is_exist bool
	is_exist, _ = ic.ipt.ChainExists("nat", KUBE_SERVICES)
	if is_exist {
		ic.ipt.ClearChain("nat", KUBE_SERVICES)
	} else {
		ic.ipt.NewChain("nat", KUBE_SERVICES)
	}
	is_exist, _ = ic.ipt.ChainExists("nat", KUBE_NODEPORTS)
	if is_exist {
		ic.ipt.ClearChain("nat", KUBE_NODEPORTS)
	} else {
		ic.ipt.NewChain("nat", KUBE_NODEPORTS)
	}

	// 在OUTPUT和PREROUTING链中插入KUBE-SERVICES,如果已经被加入则无需重复加入
	ic.ipt.AppendUnique("nat", PREROUTING, "-j", KUBE_SERVICES)
	ic.ipt.AppendUnique("nat", OUTPUT, "-j", KUBE_SERVICES)

	// // 在KUBE-SERVCES中加入KUBE_NODEPORTS
	// // 注意：KUBE_NODEPORTS必须作为KUBE-SERVCES的最后一条规则，即所有ClusterIP之后
	// ic.ipt.Append("nat", KUBE_SERVICES, "-j", KUBE_NODEPORTS)

	// 清空或创建KUBE-POSTROUTING空链并重新配置，后加入POSTROUTING,如果已经被加入则无需重复加入
	is_exist, _ = ic.ipt.ChainExists("nat", KUBE_POSTROUING)
	if is_exist {
		ic.ipt.ClearChain("nat", KUBE_POSTROUING)
	} else {
		ic.ipt.NewChain("nat", KUBE_POSTROUING)
	}
	ic.ipt.Append("nat", KUBE_POSTROUING, "-j", "MASQUERADE", "-m", "mark", "--mark", "0x4000/0x4000")
	ic.ipt.AppendUnique("nat", POSTROUING, "-j", KUBE_POSTROUING)

	// 清空或创建KUBE_MARK_MASQ空链并重新配置
	is_exist, _ = ic.ipt.ChainExists("nat", KUBE_MARK_MASQ)
	if is_exist {
		ic.ipt.ClearChain("nat", KUBE_MARK_MASQ)
	} else {
		ic.ipt.NewChain("nat", KUBE_MARK_MASQ)
	}
	ic.ipt.Append("nat", KUBE_MARK_MASQ, "-j", "MARK", "--or-mark", "0x4000")

	// 清空并删除所有的KUBE-SEP-XXX
	chain_list, _ := ic.ipt.ListChains("nat")
	for _, chain_item := range chain_list {
		if strings.HasPrefix(chain_item, KUBE_SEP) {
			ic.ipt.ClearChain("nat", chain_item)
			ic.ipt.DeleteChain("nat", chain_item)
		}
	}

	// 清空并删除所有的KUBE-SVC-XXX
	chain_list, _ = ic.ipt.ListChains("nat")
	for _, chain_item := range chain_list {
		if strings.HasPrefix(chain_item, KUBE_SVC) {
			ic.ipt.ClearChain("nat", chain_item)
			ic.ipt.DeleteChain("nat", chain_item)
		}
	}

}

func getSid(uuid string, port uint16) string {
	uuid_prefix := uuid[:8]
	return uuid_prefix + strconv.Itoa(int(port))
}

func getEid(uuid string) string {
	uuid_prefix := uuid[:8]
	return uuid_prefix
}

func uniqueAddEid(e_uuid string, euuid_list *([]string), elist *httpobject.HTTPResponse_GetAllEndpoint) string {
	if _, exist := (*elist)[e_uuid]; exist {
		can_add := true
		for _, va := range *euuid_list {
			if va == e_uuid {
				can_add = false
			}
		}
		if can_add {
			(*euuid_list) = append((*euuid_list), e_uuid)
			return add_success
		} else {
			return already_exist
		}
	} else {
		return not_found
	}
}

func getPortInfo(port uint16, port_info_list *([]service.ClusterIPPortInfo)) service.ClusterIPPortInfo {
	for _, va := range *port_info_list {
		if port == va.Port {
			return va
		}
	}
	return service.ClusterIPPortInfo{Port: 0}
}

// SyncConfig 根据请求下来的数据更新本地数据结构
//
//	@receiver ic
//	@param slist 请求下来的service list
//	@param elist 请求下来的endpoint list
func (ic *IptablesController) SyncConfig(slist *httpobject.HTTPResponse_GetAllServices, elist *httpobject.HTTPResponse_GetAllEndpoint) (kptype.KpServicesStatus, error) {
	fmt.Printf("IptablesController sync config ...\n")

	// 	// HACK: 一个简单的测试用例子，写入静态配置数据
	// 	eid1_1ID := "os1_my1"
	// 	eid1_2ID := "os1_my2"
	// 	eid2_1ID := "os2_my1"
	// 	eid2_2ID := "os2_my2"
	// 	service1ID := "1"
	// 	service2ID := "2"
	// 	service1_nodeport := 34211
	// 	service2_nodeport := 34212

	// 	ic.eid2emeta[eid1_1ID] = metadata{ip: "10.5.75.2", protocol: "tcp", port: 3000}
	// 	ic.eid2emeta[eid1_2ID] = metadata{ip: "10.5.75.3", protocol: "tcp", port: 3000}
	// 	ic.eid2emeta[eid2_1ID] = metadata{ip: "10.5.88.2", protocol: "tcp", port: 3000}
	// 	ic.eid2emeta[eid2_2ID] = metadata{ip: "10.5.88.3", protocol: "tcp", port: 3000}

	// 	ic.sid2eid[service1ID] = []string{eid1_1ID, eid2_1ID}
	// 	ic.sid2eid[service2ID] = []string{eid1_2ID, eid2_2ID}

	// 	ic.sid2smeta[service1ID] = metadata{ip: "10.100.100.1", protocol: "tcp", port: 7070}
	// 	ic.sid2smeta[service2ID] = metadata{ip: "10.100.100.2", protocol: "tcp", port: 6060}

	// ic.nodeport2sid[service1_nodeport] = service1ID
	// ic.nodeport2sid[service2_nodeport] = service2ID
	new_service_status := kptype.KpServicesStatus{ClusterIP: make(map[string]kptype.KpClusterIP), NodePort: make(map[string]kptype.KpNodePort)}

	// 创建新的数据结构，如果本轮更新成功，则替代原有的数据结构，否则保留原状
	new_sid2eid := make(map[string]([]string))
	new_nodeport2sid := make(map[uint16](string))
	new_sid2smeta := make(map[string]metadata)
	new_eid2emeta := make(map[string]metadata)

	// endpoint uuid list,在最后写入eid2emeta
	euuid_list := []string{}

	//添加clusterip所需规则
	// 从status的每一条开始访问：
	for _, clusterip := range slist.ClusterIP {
		// 如果clusterip状态不符，则不予考虑
		// if clusterip.Status.Phase != service.CLUSTERIP_ENDPOINTS_ALLOCATED &&
		// 	clusterip.Status.Phase != service.CLUSTERIP_SUCCESS {
		// 	continue
		// }
		if clusterip.Status.Phase != service.CLUSTERIP_READY {
			continue
		}

		new_kpcip := kptype.KpClusterIP{Version: clusterip.Status.Version, Vports: []uint16{}}

		for key, clusterip_port := range clusterip.Status.ServicesStatus {
			// 拿到对应的spec
			port_info := getPortInfo(key, &(clusterip.Spec.Ports))
			if port_info.Port == 0 {
				// error,这个clusterip_port不存在于spec
				continue
			}
			// port_eid_list记录本轮可以真正加入的endpoint
			port_eid_list := []string{}

			// 遍历此clusterip_port中的endpointid，检查其是否存在于拉下的endpoints中
			// 若存在，加入本轮的port_eid_list和全局euuid_list
			for _, euuid_item := range clusterip_port {
				// 加入全局euuid_list，并给出加入状态
				// 如果此euuid不存在于拉下的endpoints中，则不会加入任何组中
				if uniqueAddEid(euuid_item, &euuid_list, elist) != not_found {
					// 如果此euuid存在于拉下的endpoints中
					// 则其会被唯一加入euuid_list，并加入port_eid_list中
					port_eid_list = append(port_eid_list, getEid(euuid_item))
				}
			}

			// 创建clusterip_port对应的metadata,并在sid中添加对应eid的映射
			// 注意,只有port_eid_list中有内容，才会添加这条，这代表该clusterip有真正的endpoint可以整合进入
			if len(port_eid_list) > 0 {
				sid := getSid(clusterip.Metadata.Id, port_info.Port)
				new_sid2smeta[sid] = metadata{
					ip:       clusterip.Metadata.Ip,
					protocol: port_info.Protocol,
					port:     key,
				}
				new_sid2eid[sid] = port_eid_list

				new_kpcip.Vports = append(new_kpcip.Vports, port_info.Port)
			}
		}
		if len(new_kpcip.Vports) > 0 {
			new_service_status.ClusterIP[clusterip.Metadata.Id] = new_kpcip
		}
	}

	//添加nodeport所需规则
	for _, nodeport := range slist.NodePort {
		// if nodeport.Status.Phase != service.NODEPORT_CLUSTERIP_FINISH && nodeport.Status.Phase != service.NODEPORT_SUCCESS {
		// 	continue
		// }
		if nodeport.Status.Phase != service.NODEPORT_READY {
			continue
		}

		new_kpnp := kptype.KpNodePort{
			Version: nodeport.Status.Version,
			Nports:  []uint16{},
		}

		clusteripid := nodeport.Status.ClusterIPID
		//遍历所有ports,如果对应的sid已经被加入new数据结构，则其对应的nodeport也可被加入
		for _, nodeport_port := range nodeport.Spec.Ports {
			sid := getSid(clusteripid, nodeport_port.Port)
			if _, exist := new_sid2smeta[sid]; exist {
				new_nodeport2sid[nodeport_port.NodePort] = sid

				new_kpnp.Nports = append(new_kpnp.Nports, nodeport_port.Port)
			}
		}

		if len(new_kpnp.Nports) > 0 {
			new_service_status.NodePort[nodeport.Metadata.Id] = new_kpnp
		}

	}

	// 将endpoint存入meta数据结构(这些endpoint一定被某service用到，且存在于本轮请求下来的endpoints数据中)
	for _, euuid := range euuid_list {
		endpoint_item := (*elist)[euuid]
		new_eid2emeta[getEid(euuid)] = metadata{
			ip:       endpoint_item.PodIP,
			protocol: endpoint_item.Protocol,
			port:     endpoint_item.PodPort,
		}
	}

	// 更新无误，替代原有数据结构
	ic.sid2eid = new_sid2eid
	ic.nodeport2sid = new_nodeport2sid
	ic.sid2smeta = new_sid2smeta
	ic.eid2emeta = new_eid2emeta
	return new_service_status, nil
}

// SyncIptables 根据本地IptablesController的数据同步iptables规则
//
//	@receiver ic
func (ic *IptablesController) SyncIptables() error {
	fmt.Printf("IptablesController sync iptables ...\n")

	// 删除上一版本的内容，使其恢复到Init时的空状态
	// 删除KUBE_SERVICES中的全部条目（ClusterIP和NodePorts规则）
	ic.ipt.ClearChain("nat", KUBE_SERVICES)

	// 删除KUBE_NODEPORT中的全部条目（导向KUBE_SVC_XXX）
	ic.ipt.ClearChain("nat", KUBE_NODEPORTS)

	// 删除所有KUBE_SEP_XXX链中的所有条目，并删除此链
	nat_chain_list, _ := ic.ipt.ListChains("nat")
	for _, chain := range nat_chain_list {
		if strings.HasPrefix(chain, KUBE_SEP) {
			ic.ipt.ClearChain("nat", chain)
			ic.ipt.DeleteChain("nat", chain)
		}
	}

	// 删除所有KUBE_SVC_XXX链中的所有条目，并删除此链
	nat_chain_list, _ = ic.ipt.ListChains("nat")
	for _, chain := range nat_chain_list {
		if strings.HasPrefix(chain, KUBE_SVC) {
			ic.ipt.ClearChain("nat", chain)
			ic.ipt.DeleteChain("nat", chain)
		}
	}

	// 根据本地map中的数据，重新构建iptables规则
	// 创建所有KUBE_SEP_XXX链
	for key, value := range ic.eid2emeta {
		ic.ipt.NewChain("nat", KUBE_SEP+key)
		ic.ipt.Append("nat", KUBE_SEP+key, "-m", "addrtype", "--src-type", "LOCAL", "-j", KUBE_MARK_MASQ)
		ic.ipt.Append("nat", KUBE_SEP+key, "-s", value.ip, "-j", KUBE_MARK_MASQ)
		ic.ipt.Append("nat", KUBE_SEP+key, "-p", value.protocol, "-m", value.protocol, "-j", DNAT, "--to-destination", value.ip+":"+strconv.Itoa(int(value.port)))
	}

	// 创建所有KUBE_SVC_XXX链
	for sid, eid_list := range ic.sid2eid {
		ic.ipt.NewChain("nat", KUBE_SVC+sid)
		// 填入endpoints规则
		// 注意：第一条规则需要指明负载均衡配置
		length := len(eid_list)
		for index, eid_item := range eid_list {
			var pro float32 = 1.0 / float32(length-index)
			pro_string := fmt.Sprintf("%.2f", pro)
			ic.ipt.Append("nat", KUBE_SVC+sid, "-m", "statistic", "--mode", "random", "--probability", pro_string, "-j", KUBE_SEP+eid_item)
		}
	}

	// 在KUBE-NODEPORTS中填入MARK-SNAT规则和service转发规则
	for port, sid := range ic.nodeport2sid {
		meta := (ic.sid2smeta)[sid]
		ic.ipt.Append("nat", KUBE_NODEPORTS, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(int(port)), "-j", KUBE_MARK_MASQ)
		ic.ipt.Append("nat", KUBE_NODEPORTS, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(int(port)), "-j", KUBE_SVC+sid)
	}

	// 在KUBE-SERVICES中填入规则，先ClusterIP后NodePort
	for sid, meta := range ic.sid2smeta {
		ic.ipt.Append("nat", KUBE_SERVICES, "-d", meta.ip, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(int(meta.port)), "-j", KUBE_SVC+sid)
	}
	ic.ipt.Append("nat", KUBE_SERVICES, "-m", "addrtype", "--dst-type", "LOCAL", "-j", KUBE_NODEPORTS)

	// 此条规则需要解决的问题已经KUBE-SEP-XXX的第一条规则上
	// 若按原配置，访问部分公网服务（如canvas等）会无法连接
	// // 在service向外机负载均衡时，配置SNAT
	// ic.ipt.Insert("nat", POSTROUING, 1, "-m", "addrtype", "--src-type", "LOCAL", "-j", "MASQUERADE")
	return nil
}
