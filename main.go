package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) fileserverHits_checkup(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hits))
}

func (cfg *apiConfig) fileserverHits_Reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server Hits Reset to 0"))
}

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	apiCfg := new(apiConfig)

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc((http.StripPrefix("/app/", http.FileServer(http.Dir("."))))))
	serveMux.Handle("assets/logo.png", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	serveMux.HandleFunc("/metrics", apiCfg.fileserverHits_checkup)
	serveMux.HandleFunc("/reset", apiCfg.fileserverHits_Reset)
	http.ListenAndServe(server.Addr, server.Handler)
}
