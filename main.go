package main

import (
	"fmt"
	"io"
	"net/http"
)

func proxyhandler(w http.ResponseWriter, r *http.Request) {
	//copy the incoming request and change the host to backend
	clone := r.Clone(r.Context())
	clone.URL.Scheme = "http"
	clone.URL.Host = "localhost:8081"
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
	defer res.Body.Close()
}
func main() {
	mux := http.NewServeMux()
	fmt.Print("starting server\n")
	mux.HandleFunc("/", proxyhandler)
	http.ListenAndServe(":3000", mux)
}
