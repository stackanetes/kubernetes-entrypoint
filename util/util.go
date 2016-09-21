package util

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetIp() (ip string, err error) {
	iface := os.Getenv("INTERFACE_NAME")
	if iface == "" {
		return "", fmt.Errorf("Environment variable INTERFACE_NAME not set")
	}
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return "", fmt.Errorf("Cannot get iface: %v", err)
	}

	address, err := i.Addrs()
	if err != nil || len(address) == 0 {
		return "", fmt.Errorf("Cannot get ip: %v", err)
	}
	//Take first element to get rid of subnet
	ip = strings.Split(address[0].String(), "/")[0]
	return
}
