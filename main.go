package main

import (
	"net/http"
)

func main() {

	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	serveMux.Handle("assets/logo.png", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(server.Addr, server.Handler)

}
