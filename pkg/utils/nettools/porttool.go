package nettools

import (
	"fmt"
	"net"
)

// CheckPortAvailability try to listen on the given port to check if it's available
func CheckPortAvailability(port int) bool {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func GetAvailablePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	return addr.Port, nil
}
