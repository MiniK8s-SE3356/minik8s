package nginxManager

import (
	"fmt"
	"os"
	"os/exec"

	httpobject "github.com/MiniK8s-SE3356/minik8s/pkg/types/httpObject"
)

var nginx_defconf_file string = "/etc/nginx/sites-available/default"

type NginxManager struct {
}

func NewDnsManager() *NginxManager {
	fmt.Printf("New NewDnsManager\n")
	return &NginxManager{}
}

func (nm *NginxManager) Init() {
	fmt.Printf("Init NginxManager\n")
}

func (nm *NginxManager) SyncNginx(dns_list *httpobject.HTTPResponse_GetAllDns) {
	fmt.Printf("NginxManager SyncNginx\n")

	// fmt.Println((*dns_list))

	whole_dns_server := ""

	for _, dns_item := range *dns_list {
		path_proxy_list := ""
		for _, path_status_item := range dns_item.Status.PathsStatus {
			path_proxy_item := fmt.Sprintf(custom_path_proxy, path_status_item.SubPath, path_status_item.SvcIP, int(path_status_item.SvcPort))
			// fmt.Println(path_proxy_item)
			path_proxy_list += path_proxy_item
			// fmt.Println(path_proxy_list)
		}

		dns_server := fmt.Sprintf(custom_server, dns_item.Spec.Host, path_proxy_list)
		fmt.Println(dns_server)
		whole_dns_server = whole_dns_server + dns_server
	}

	whole_content := fmt.Sprintf(nginx_default, whole_dns_server)
	// fmt.Printf(whole_content)

	// 字符串组装完成，写入对应的配置文件中
	file, err := os.OpenFile(nginx_defconf_file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Filed to open the nginx file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(whole_content)
	if err != nil {
		fmt.Println("Filed to write content to nginx file:", err)
		return
	}

	// TODO: 更新配置后，需要使用nginx -t检验配置无误，如果出现错误，不应该重启nginx,以防止nginx crash

	// 更新完成之后，重启nginx系统进程，以刷新配置
	cmd := exec.Command("systemctl", "restart", "nginx")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("restart nginx failed:", err)
	}
}
