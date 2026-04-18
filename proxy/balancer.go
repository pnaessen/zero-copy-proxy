package proxy

import (
	"sync/atomic"
)

type RoundRobin struct {
	targets []string
	counter uint64
}

func NewRoundRobin(targets []string) *RoundRobin {
	return &RoundRobin{
		targets: targets,
		counter: 0,
	}
}

func (rr *RoundRobin) NextTarget() string {

	val := atomic.AddUint64(&rr.counter, 1)

	index := (val - 1) % uint64(len(rr.targets))

	return rr.targets[index]
}