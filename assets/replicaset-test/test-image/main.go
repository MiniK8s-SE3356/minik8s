package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			http.Error(w, "Could not get IP address", http.StatusInternalServerError)
			return
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Fprintln(w, ipnet.IP.String())
				}
			}
		}
	})

	fmt.Println("Server is running on port 10086")
	err := http.ListenAndServe(":10086", nil)
	if err != nil {
		fmt.Println("Failed to start server: ", err)
		return
	}
}
