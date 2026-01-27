package httpx

import (
	"encoding/json"
	"net/http"

	"crypto-api/internal/models"
)

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Error: message,
	}
	json.NewEncoder(w).Encode(resp)
}
