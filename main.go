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
	metricsAddr := ":8081"

	serveursCibles := []string{
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	}

	lb := proxy.NewRoundRobin(serveursCibles)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		reqCount, activeCount := lb.GetStats()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		fmt.Fprintf(w, "# HELP proxy_requests_total Total number of handled requests\n")
		fmt.Fprintf(w, "# TYPE proxy_requests_total counter\n")
		fmt.Fprintf(w, "proxy_requests_total %d\n", reqCount)

		fmt.Fprintf(w, "# HELP proxy_active_servers Number of healthy backend servers\n")
		fmt.Fprintf(w, "# TYPE proxy_active_servers gauge\n")
		fmt.Fprintf(w, "proxy_active_servers %d\n", activeCount)
	})

	go func() {
		log.Printf("Prometheus metrics available at http://localhost%s/metrics", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Failed to start TCP listener on %s: %v", localAddr, err)
	}
	defer listener.Close()

	fmt.Printf("TCP proxy started on %s\n", localAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept client connection: %v", err)
			continue
		}

		cibleActuelle := lb.NextTarget()
		if cibleActuelle == "" {
			log.Printf("No backend available, closing client connection")
			clientConn.Close()
			continue
		}

		go proxy.HandleConnection(clientConn, cibleActuelle)
	}
}
