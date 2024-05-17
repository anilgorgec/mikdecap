package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anilgorgec/mikdecap/internal/adap"
	"github.com/anilgorgec/mikdecap/internal/tzsp"
)

const MAX_WORKER = 2

func openListener(listenAddr string) *net.UDPConn {

	addr, _ := net.ResolveUDPAddr("udp", listenAddr)
	conn, _ := net.ListenUDP("udp", addr)
	return conn

}
func worker(fd int, addrLink syscall.Sockaddr, wg *sync.WaitGroup, jobs <-chan []byte) {
	defer wg.Done()
	buf := make([]byte, 65535)
	for j := range jobs {
		ln, err := tzsp.Parse(j, buf)
		if err != nil {
			fmt.Println("Parse Error : " + hex.Dump(j))
			continue
		}

		if err := syscall.Sendto(fd, buf[:ln], 0, addrLink); err != nil {
			fmt.Printf("Full HexDump : Len(%d) - %s", len(j), hex.Dump(j))
			fmt.Println(err)
		}
	}
}

func reader(ctx context.Context, conn *net.UDPConn, jobs chan<- []byte) {
	buf := make([]byte, 65535)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break loop
				}
				fmt.Println(err)
				os.Exit(1)
			}
			jobs <- buf[:n]

		}
	}
	fmt.Println("Consumer received cancellation signal, closing channel!")
	close(jobs)
	fmt.Println("Consumer closed channel")
}
func main() {
	listenAddr := flag.String("l", ":37008", "Address to listen for TZSP UDP packets.")
	iface := flag.String("i", "dummy0", "destination interface")
	flag.Parse()
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	fd, addrLink := adap.OpenRawSocket(*iface)
	defer syscall.Close(fd)
	conn := openListener(*listenAddr)
	jobs := make(chan []byte)

	for w := 0; w < MAX_WORKER; w++ {
		wg.Add(1)
		go worker(fd, addrLink, wg, jobs)
	}
	go reader(ctx, conn, jobs)

	fmt.Printf("Listener : %s, Interface: %s\n", *listenAddr, *iface)
	<-termChan
	_ = conn.Close()
	cancelFunc()
	wg.Wait()
	fmt.Println("Workers done, shutting down!")

}
