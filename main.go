package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pnaessen/zero-copy-proxy/proxy"
)

func main() {
	localAddr := ":8080"
	targetAddr := "google.com:80"

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

		go proxy.HandleConnection(clientConn, targetAddr)
	}
}
