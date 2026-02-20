package httpx

import (
	"crypto-api/internal/storage"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func HealthHandler(store *storage.PostgresStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := "OK"
		statusCode := http.StatusOK

		// Protect against nil store or pool
		if store == nil || store.Pool == nil {
			status = "Database Configuration Missing"
			statusCode = http.StatusServiceUnavailable
		} else {
			// Perform real ping
			if err := store.Pool.Ping(r.Context()); err != nil {
				status = "Database Disconnected"
				statusCode = http.StatusServiceUnavailable
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": status,
			"time":   time.Now().Format(time.RFC3339),
		}); err != nil {
			log.Printf("[http] failed to encode health response: %v", err)
		}
	}
}
