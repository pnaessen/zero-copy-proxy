package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pnaessen/zero-copy-proxy/proxy"
)

func main() {
	localAddr := ":8123"
	targetAddr := []string{
		"httpbin.org:80",
		"google.com:80",
		"truc.com:80",
	}

	lb := proxy.NewRoundRobin(targetAddr)

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Failed to start TCP listener: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Zero-copy proxy started on %s (forwarding to %s)\n", localAddr, targetAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}

		target := lb.NextTarget()

		go proxy.HandleConnection(clientConn, target)
	}
}
