package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jb-8200/cursor-analytics-platform/services/cursor-sim/internal/config"
)

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func main() {
	// Parse configuration from flags
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Log configuration
	log.Printf("cursor-sim starting with configuration:")
	log.Printf("  Port: %d", cfg.Port)
	log.Printf("  Developers: %d", cfg.Developers)
	log.Printf("  Velocity: %s", cfg.Velocity)
	log.Printf("  Fluctuation: %.2f", cfg.Fluctuation)
	log.Printf("  Seed: %d", cfg.Seed)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/v1/health", handleHealth)

	// Placeholder endpoints - to be implemented
	mux.HandleFunc("/v1/analytics/ai-code/commits", handleNotImplemented)
	mux.HandleFunc("/v1/analytics/ai-code/changes", handleNotImplemented)
	mux.HandleFunc("/v1/analytics/team/agent-edits", handleNotImplemented)
	mux.HandleFunc("/v1/analytics/team/tabs", handleNotImplemented)
	mux.HandleFunc("/v1/analytics/team/dau", handleNotImplemented)
	mux.HandleFunc("/v1/analytics/team/models", handleNotImplemented)
	mux.HandleFunc("/v1/teams/members", handleNotImplemented)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("cursor-sim listening on %s", addr)
	log.Printf("Health check: http://localhost:%d/v1/health", cfg.Port)
	log.Printf("Status: TASK-SIM-002 complete - CLI flag parsing implemented")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Service:   "cursor-sim",
		Version:   "0.0.1-p0",
	}

	writeJSON(w, http.StatusOK, response)
}

func handleNotImplemented(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED",
		fmt.Sprintf("Endpoint %s not yet implemented (P0 scaffolding)", r.URL.Path))
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	writeJSON(w, status, response)
}
