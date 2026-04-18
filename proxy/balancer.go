package proxy

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type RoundRobin struct {
	allTargets    []string
	activeTargets []string
	mu            sync.RWMutex
	counter       uint64
}

func NewRoundRobin(targets []string) *RoundRobin {
	rr := &RoundRobin{
		allTargets:    targets,
		activeTargets: targets,
		counter:       0,
	}

	go rr.healthCheckLoop()

	return rr
}

func (rr *RoundRobin) NextTarget() string {
	rr.mu.RLock()
	defer rr.mu.RUnlock()

	if len(rr.activeTargets) == 0 {
		return ""
	}

	val := atomic.AddUint64(&rr.counter, 1)
	index := (val - 1) % uint64(len(rr.activeTargets))

	return rr.activeTargets[index]
}

func (rr *RoundRobin) healthCheckLoop() {
	for {
		var healthy []string

		for _, target := range rr.allTargets {
			conn, err := net.DialTimeout("tcp", target, 2*time.Second)
			if err != nil {
				log.Printf("[HealthCheck] Serveur dead : %s", target)
			} else {
				conn.Close()
				healthy = append(healthy, target)
			}
		}

		rr.mu.Lock()
		rr.activeTargets = healthy
		rr.mu.Unlock()

		time.Sleep(5 * time.Second)
	}
}

func (rr *RoundRobin) GetStats() (uint64, int) {
	totalRequests := atomic.LoadUint64(&rr.counter)

	rr.mu.RLock()
	activeCount := len(rr.activeTargets)
	rr.mu.RUnlock()

	return totalRequests, activeCount
}