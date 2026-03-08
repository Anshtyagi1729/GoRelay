package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var backends = []string{
	"localhost:8081",
	"localhost:8082",
	"localhost:8083",
}
var health = map[string]bool{
	"localhost:8081": true,
	"localhost:8082": true,
	"localhost:8083": true,
}
var current atomic.Int64
var healthMu sync.RWMutex

func healthchecks() {
	for {
		time.Sleep(10 * time.Second)
		for _, backend := range backends {
			_, err := http.Get("http://" + backend)
			if err != nil {
				healthMu.Lock()
				health[backend] = false
				healthMu.Unlock()
			} else {
				healthMu.Lock()
				health[backend] = true
				healthMu.Unlock()
			}
		}
	}
}
func proxyhandler(w http.ResponseWriter, r *http.Request) {
	//copy the incoming request and change the host to backend

	clone := r.Clone(r.Context())
	idx := current.Add(1) % int64(len(backends))
	host := backends[idx]
	healthMu.RLock()
	ishealthy := health[host]
	healthMu.RUnlock()
	tried := 0
	for !ishealthy {
		if tried == len(backends) {
			http.Error(w, "no health backend", http.StatusServiceUnavailable)
			return
		}
		idx = current.Add(1) % int64(len(backends))
		host = backends[idx]
		healthMu.RLock()
		ishealthy = health[host]
		healthMu.RUnlock()
		tried++
	}
	clone.URL.Scheme = "http"
	clone.URL.Host = backends[idx]
	//roundtrip and handle the error
	res, err := http.DefaultTransport.RoundTrip(clone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	for k, vv := range res.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
	res.Body.Close()
}
func main() {
	mux := http.NewServeMux()
	fmt.Print("starting proxy server at :3000\n")
	go healthchecks()
	mux.HandleFunc("/", proxyhandler)
	http.ListenAndServe(":3000", mux)
}
