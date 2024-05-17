//go:build linux

package adap

import (
	"fmt"
	"net"
	"syscall"
)

func OpenRawSocket(iface string) (int, syscall.Sockaddr) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.IPPROTO_RAW)))
	if err != nil {
		panic(fmt.Errorf("failed to create socket: %w", err))
	}

	// err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 524288)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to create socket: %w", err))
	// }
	ifi, err := net.InterfaceByName(iface)
	if err != nil {
		panic(fmt.Errorf("failed to get interface %s: %w", iface, err))
	}

	return fd, &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.IPPROTO_RAW),
		Ifindex:  ifi.Index,
	}
}
