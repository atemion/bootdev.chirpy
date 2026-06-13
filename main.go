package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	// mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", myHandler)
	mux.HandleFunc("GET /metrics", apiCfg.serverHitsHandler)
	mux.HandleFunc("POST /reset", apiCfg.serverHitsReset)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	s.ListenAndServe()
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Middleware to log server hits

// Define struct to store value
type apiConfig struct {
	fileserverHits atomic.Int32
}

// wrap main handler in middleware to increment with every hit
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// handler to fetch server hits via /metrics
func (cfg *apiConfig) serverHitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	numHits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(numHits))
}

// handler to reset server hits via /reset
func (cfg *apiConfig) serverHitsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Server hit count reset!"))
}
