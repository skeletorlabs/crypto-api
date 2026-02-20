package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"crypto-api/internal/models"
)

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[http] failed to encode error response (status=%d): %v", status, err)
	}
}
