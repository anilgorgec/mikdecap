//go:build darwin

package adap

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
)

func transformIPArray(ip net.Addr) [4]byte {
	var s [4]byte
	ipArr := strings.Split(ip.String(), "/")
	ipArr = strings.Split(ipArr[0], ".")
	for i, pi := range ipArr {
		v, _ := strconv.Atoi(pi)
		s[i] = byte(v)
	}
	return s
}

func OpenRawSocket(iface string) (int, syscall.Sockaddr) {

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		panic(fmt.Errorf("failed to create socket: %w", err))
	}
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		panic(err)
	}
	ifi, err := net.InterfaceByName(iface)

	if err != nil {
		panic(fmt.Errorf("failed to get interface %s: %w", iface, err))
	}
	addrs, _ := ifi.Addrs()

	return fd, &syscall.SockaddrInet4{Addr: transformIPArray(addrs[0])}
}
