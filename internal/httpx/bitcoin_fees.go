package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/bitcoin"
	"encoding/json"
	"net/http"
	"time"
)

func BitcoinFeesHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check cache first
		if cached, ok := c.Get("fees"); ok {
			resp := cached.(models.BitcoinFeesResponse)
			resp.Cached = true
			json.NewEncoder(w).Encode(resp)
			return
		}

		ctx := r.Context()
		fees, err := bitcoin.GetBitcoinFees(ctx)
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		resp := models.BitcoinFeesResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			Low:    fees.HourFee,
			Medium: fees.HalfHourFee,
			High:   fees.FastestFee,
		}

		c.Set("fees", resp, 30*time.Second)
		json.NewEncoder(w).Encode(resp)
	}
}
