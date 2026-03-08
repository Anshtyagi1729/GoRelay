package main

import (
	"log"
	"net/http"
	"time"
)

func healthchecks() {
	for {
		time.Sleep(10 * time.Second)
		for _, backend := range backends {
			_, err := http.Get("http://" + backend.Address)
			if err != nil {
				backendMu.Lock()
				backend.Healthy = false
				backendMu.Unlock()
				log.Printf("backend %s is down", backend.Address)
			} else {
				backendMu.Lock()
				backend.Healthy = true
				backendMu.Unlock()
			}
		}
	}
}
