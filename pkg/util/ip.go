package util

import (
	"net"
	"strings"
)

const (
	Subnetwork_192_168 = "192.168."
	Subnetwork_10      = "10."
	Subnetwork_lo      // 本地回环子网 127.0.0.1
)

func GetLocalIp() (ip string) {
	if len(ip) == 0 {
		ip = GetIpv4_192_168()
	}
	if len(ip) == 0 {
		ip = GetIpv4_172()
	}
	if len(ip) == 0 {
		ip = GetIpv4_10()
	}
	return
}

func GetIpv4_192_168() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), Subnetwork_192_168) == 0 {
				return sip
			}
		}
	}
	return ""
}
func GetIpv4_172() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "172.") == 0 {
				return sip
			}
		}
	}
	return ""
}
func GetIpv4_10() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "10.") == 0 {
				return sip
			}
		}
	}
	return ""
}
func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}
	return ips
}
