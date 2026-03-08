package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	Address string
	Healthy bool
}
type ClientInfo struct {
	Count       int
	WindowStart time.Time
}

var clients = map[string]*ClientInfo{}
var clientMu sync.RWMutex

var backends = []*Backend{
	{Address: "localhost:8081", Healthy: true},
	{Address: "localhost:8082", Healthy: true},
	{Address: "localhost:8083", Healthy: true},
}
var backendMu sync.RWMutex
var current atomic.Int64

func rateLimiter(w http.ResponseWriter, r *http.Request) bool {
	client, _, _ := net.SplitHostPort(r.RemoteAddr)
	// fmt.Printf("client:%s\n", client)
	clientMu.RLock()
	_, exists := clients[client]
	clientMu.RUnlock()
	if !exists {
		clientMu.Lock()
		clients[client] = &ClientInfo{Count: 1, WindowStart: time.Now()}
		clientMu.Unlock()
		return true
	} else {
		clientMu.Lock()
		if time.Since(clients[client].WindowStart) > 10*time.Second {
			clients[client].Count = 1
			clients[client].WindowStart = time.Now()
		} else if clients[client].Count > 100 {
			clientMu.Unlock()
			http.Error(w, "Rate limited", http.StatusTooManyRequests)
			return false
		} else {
			clients[client].Count++
			// fmt.Printf("clients:%s,count:%d", client, clients[client].Count)
		}
		clientMu.Unlock()
		return true
	}
}
func healthchecks() {
	for {
		time.Sleep(10 * time.Second)
		for _, backend := range backends {
			_, err := http.Get("http://" + backend.Address)
			if err != nil {
				backendMu.Lock()
				backend.Healthy = false
				backendMu.Unlock()
			} else {
				backendMu.Lock()
				backend.Healthy = true
				backendMu.Unlock()
			}
		}
	}
}
func proxyhandler(w http.ResponseWriter, r *http.Request) {
	//copy the incoming request and change the host to backend
	if !rateLimiter(w, r) {
		return
	}
	clone := r.Clone(r.Context())
	idx := current.Add(1) % int64(len(backends))
	host := backends[idx]
	backendMu.RLock()
	ishealthy := host.Healthy
	backendMu.RUnlock()
	tried := 0
	for !ishealthy {
		if tried == len(backends) {
			http.Error(w, "no health backend", http.StatusServiceUnavailable)
			return
		}
		idx = current.Add(1) % int64(len(backends))
		host = backends[idx]
		backendMu.RLock()
		ishealthy = host.Healthy
		backendMu.RUnlock()
		tried++
	}
	clone.URL.Scheme = "http"
	clone.URL.Host = backends[idx].Address
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
