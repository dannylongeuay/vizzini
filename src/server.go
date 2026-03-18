package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	SERVER_PORT            = ":8080"
	DEFAULT_ALLOWED_ORIGIN = "https://chess.cyberdan.dev"
)

func ModeServe() {
	allowedOrigin := DEFAULT_ALLOWED_ORIGIN
	if os.Getenv("CORS_PERMISSIVE") != "" {
		allowedOrigin = "*"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)

	handler := corsMiddleware(mux, allowedOrigin)

	fmt.Printf("Listening on %s (CORS origin: %s)\n", SERVER_PORT, allowedOrigin)
	if err := http.ListenAndServe(SERVER_PORT, handler); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func corsMiddleware(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
