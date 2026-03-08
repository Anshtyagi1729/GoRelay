package main

import (
	"net/http"
)

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
	mux1.HandleFunc("/", backend1hand)
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/", backend2hand)
	go http.ListenAndServe(":8082", mux1)
	go http.ListenAndServe(":8083", mux2)
	http.ListenAndServe(":8081", nil)
}

//TODO:
// rate limiting
// load balancing
// caching
