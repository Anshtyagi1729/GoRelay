package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	Address  string
	Healthy  bool
	Requests atomic.Int64
	Latency  time.Duration
}

var backends []*Backend
var backendMu sync.RWMutex

func addBackend(address string) {
	for _, b := range backends {
		if b.Address == address {
			return
		}
	}
	backends = append(backends, &Backend{Address: address, Healthy: true})
}

func listbackend() []*Backend {
	backendMu.RLock()
	defer backendMu.RUnlock()
	return backends
}
