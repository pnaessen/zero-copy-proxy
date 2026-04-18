package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/pnaessen/zero-copy-proxy/proxy"
)

func main() {
	localAddr := ":8123"
	metricsAddr := ":8124"

	targetAddr := []string{
		"httpbin.org:80",
		"google.com:80",
		"truc.com:80",
	}

	lb := proxy.NewRoundRobin(targetAddr)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		reqCount, activeCount := lb.GetStats()
		fmt.Fprintf(w, "Load Balancer Zero-Copy\n")
		fmt.Fprintf(w, "-----------------------\n")
		fmt.Fprintf(w, "Requêtes traitées : %d\n", reqCount)
		fmt.Fprintf(w, "Serveurs actifs   : %d / %d\n", activeCount, len(targetAddr))
	})

	go func() {
		log.Printf(" Serveur on http://localhost%s/metrics\n", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Fatalf("Erreur serveur métriques : %v", err)
		}
	}()

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

		if target == "" {
			log.Printf("No serveur up close client: %s", clientConn.RemoteAddr())
			clientConn.Close()
			continue
		}

		go proxy.HandleConnection(clientConn, target)
	}
}
