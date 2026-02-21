package bitcoin

import (
	api "crypto-api/internal/http"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"log"
	"net/http"
)

func BitcoinCorrelationHandler(
	repo *repositories.IntelligenceRepository,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		snapshot, err := repo.GetLatest(r.Context())
		if err != nil || snapshot == nil {
			api.JSONError(w, http.StatusServiceUnavailable, "Correlation data unavailable")
			return
		}

		response := struct {
			Meta models.Meta `json:"meta"`
			Data struct {
				Correlation float64 `json:"correlation"`
				SampleDate  string  `json:"snapshotDate"`
			} `json:"data"`
		}{
			Meta: models.Meta{
				UpdatedAt: snapshot.CreatedAt,
				Cached:    false,
			},
		}

		response.Data.Correlation = snapshot.Correlation
		response.Data.SampleDate = snapshot.SnapshotDate.Format("2006-01-02")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("[http] failed to encode correlation response: %v", err)
		}
	}
}
