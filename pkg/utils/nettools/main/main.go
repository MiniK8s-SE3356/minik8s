package main

import (
	"fmt"
	"net"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		for _, addr := range addrs {
			// check address type and get the IP address
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// filter out loopback and non-ipv4 addresses
			if ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				fmt.Printf("Interface: %s, IP: %s\n", i.Name, ip.String())
			}
		}
	}
}
