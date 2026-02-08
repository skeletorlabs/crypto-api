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
		cacheKey := cache.KeyBitcoinFees

		// Check cache first
		if cached, ok := cache.Get[models.BitcoinFeesResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			json.NewEncoder(w).Encode(cached)
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

		cache.Set(c, cacheKey, resp, 30*time.Second)
		json.NewEncoder(w).Encode(resp)
	}
}
