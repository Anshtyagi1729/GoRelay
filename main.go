package main

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type ClientInfo struct {
	Count       int
	WindowStart time.Time
}

var clients = map[string]*ClientInfo{}
var clientMu sync.RWMutex

var current atomic.Int64

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config:%v\n", err)
	}
	for _, addr := range cfg.Backends {
		backends = append(backends, &Backend{Address: addr, Healthy: true})
	}
	mux := http.NewServeMux()
	log.Printf("reverse proxy started at %s for %d backend", cfg.Port, len(backends))
	go healthchecks()
	mux.HandleFunc("/", proxyhandler)
	http.ListenAndServe(cfg.Port, mux)
}
