package utils

import (
	"net"
)

func GetLocalIP() string {
	netInterfaces, err := net.Interfaces()
	if HasErr(err) {
		return ""
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if !IsEmpty(ipnet.IP.To4()) {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}
