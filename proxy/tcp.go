package proxy

import (
	"io"
	"log"
	"net"
	"sync"
)

func HandleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	clientConn.SetReadDeadline(time.Now().Add(30 * time.Second))
	
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v", targetAddr, err)
		return
	}
	defer targetConn.Close()

	log.Printf("New proxied connection: %s <-> %s", clientConn.RemoteAddr(), targetAddr)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := io.Copy(targetConn, clientConn)
		if err != nil && err != io.EOF {
			log.Printf("Error forwarding client to target: %v", err)
		}
		targetConn.(*net.TCPConn).CloseWrite()
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(clientConn, targetConn)
		if err != nil && err != io.EOF {
			log.Printf("Error forwarding target to client: %v", err)
		}
		clientConn.(*net.TCPConn).CloseWrite()
	}()

	wg.Wait()
	log.Printf("Connection closed: %s", clientConn.RemoteAddr())
}
