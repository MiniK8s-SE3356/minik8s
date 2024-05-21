package nettools

import (
	"fmt"
	"net"
)

func GetAllInterface() (map[string]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	interfaceMap := make(map[string]string)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				interfaceMap[i.Name] = ip.String()
			}
		}
	}
	return interfaceMap, nil
}

func GetLocalIP(interface_name string) (string, error) {

	mapInterface, err := GetAllInterface()
	if err != nil {
		return "", err
	}

	if ip, ok := mapInterface[interface_name]; ok {
		return ip, nil
	}

	return "", fmt.Errorf("interface %s not found", interface_name)
}

func KubeletDefaultIP() string {
	ip, _ := GetLocalIP("ens3")
	if ip != "" {
		return ip
	}

	interfaces, err := GetAllInterface()
	if err != nil {
		return "127.0.0.1"
	}
	if len(interfaces) == 0 {
		return "127.0.0.1"
	}

	for _, v := range interfaces {
		return v
	}

	return "127.0.0.1"
}
