package iptablesController

import (
	"fmt"
	"strconv"
	"strings"

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
)

type metadata struct {
	ip       string
	protocol string /* 只允许tcp和udp两种字符串 */
	port     int
}

type IptablesController struct {
	// 控制iptables的实际对象
	ipt *iptables.IPTables

	// 记录service下辖pods，service ID --> endpoint ID数组
	// ID是该类对象的集群唯一标识符，这些符号会参与构建KUBE-SVC-XXX和KUBE-SEP-XXX
	sid2eid map[string]([]string)

	// 记录nodePort与service对应关系，本机端口号 --> service ID
	nodeport2sid map[int](string)

	// service ID --> service metadata的map
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
	ic.nodeport2sid = make(map[int]string)
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

// SyncConfig 根据传入的参数更新本地IptablesController数据
//
//	@receiver ic
func (ic *IptablesController) SyncConfig() {
	// TODO: 设计并完成SyncConfig
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
}

// SyncIptables 根据本地IptablesController的数据同步iptables规则
//
//	@receiver ic
func (ic *IptablesController) SyncIptables() {
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
		ic.ipt.Append("nat", KUBE_SEP+key, "-p", value.protocol, "-m", value.protocol, "-j", DNAT, "--to-destination", value.ip+":"+strconv.Itoa(value.port))
	}

	// 创建所有KUBE_SVC_XXX链
	for sid, eid_list := range ic.sid2eid {
		ic.ipt.NewChain("nat", KUBE_SVC+sid)
		// 填入endpoints规则
		// 注意：第一条规则需要指明负载均衡配置
		for index, eid_item := range eid_list {
			if index == 0 {
				var pro float32 = 1.0 / float32(len(eid_list))
				pro_string := fmt.Sprintf("%.2f", pro)
				ic.ipt.Append("nat", KUBE_SVC+sid, "-m", "statistic", "--mode", "random", "--probability", pro_string, "-j", KUBE_SEP+eid_item)
			} else {
				ic.ipt.Append("nat", KUBE_SVC+sid, "-j", KUBE_SEP+eid_item)
			}
		}
	}

	// 在KUBE-NODEPORTS中填入MARK-SNAT规则和service转发规则
	for port, sid := range ic.nodeport2sid {
		meta := (ic.sid2smeta)[sid]
		ic.ipt.Append("nat", KUBE_NODEPORTS, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(port), "-j", KUBE_MARK_MASQ)
		ic.ipt.Append("nat", KUBE_NODEPORTS, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(port), "-j", KUBE_SVC+sid)
	}

	// 在KUBE-SERVICES中填入规则，先ClusterIP后NodePort
	for sid, meta := range ic.sid2smeta {
		ic.ipt.Append("nat", KUBE_SERVICES, "-d", meta.ip, "-p", meta.protocol, "-m", meta.protocol, "--dport", strconv.Itoa(meta.port), "-j", KUBE_SVC+sid)
	}
	ic.ipt.Append("nat", KUBE_SERVICES, "-m", "addrtype", "--dst-type", "LOCAL", "-j", KUBE_NODEPORTS)

	// 此条规则需要解决的问题已经KUBE-SEP-XXX的第一条规则上
	// 若按原配置，访问部分公网服务（如canvas等）会无法连接
	// // 在service向外机负载均衡时，配置SNAT
	// ic.ipt.Insert("nat", POSTROUING, 1, "-m", "addrtype", "--src-type", "LOCAL", "-j", "MASQUERADE")

}
