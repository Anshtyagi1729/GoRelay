package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ClientInfo struct {
	Count       int
	WindowStart time.Time
}

var clients = map[string]*ClientInfo{}
var clientMu sync.RWMutex

var current atomic.Int64

func main() {
	logFile, _ := os.OpenFile("gorelay.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(logFile)
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
	ti := textinput.New()
	ti.Placeholder = "localhost:8084"
	m := model{backends: listBackends(), input: ti}
	p := tea.NewProgram(m)
	go p.Run()
	http.ListenAndServe(cfg.Port, mux)
}
