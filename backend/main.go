package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><h1>hello from backend</h1></body></html>"))
	})
	http.ListenAndServe(":8081", nil)
}

//TODO:
// multiple backends
// roundrobbin
// rate limiting
// load balancing
// caching
