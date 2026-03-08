package main

import (
	"fmt"
	"net/http"
)
//only for test
// dummy file for multiple backends
func backend1hand(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><h1>hello from backend1</h1></body></html>"))
}
func backend2hand(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><h1>hello from backend2</h1></body></html>"))
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><h1>hello from backend</h1></body></html>"))
	})
	mux1 := http.NewServeMux()
	fmt.Printf("starting all the servers")
	mux1.HandleFunc("/", backend1hand)
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/", backend2hand)
	go http.ListenAndServe(":8082", mux1)
	go http.ListenAndServe(":8083", mux2)
	http.ListenAndServe(":8081", nil)
}

//TODO:
// caching
//the dashboard and further refactor
// adding other loadbalancing tricks
// adding the tls option
//
