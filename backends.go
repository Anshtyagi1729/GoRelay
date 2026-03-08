package main

import (
	"sync"
)

type Backend struct {
	Address string
	Healthy bool
}

var backends []*Backend
var backendMu sync.RWMutex
