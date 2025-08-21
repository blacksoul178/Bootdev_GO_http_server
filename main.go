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

	http.ListenAndServe(server.Addr, server.Handler)

}
