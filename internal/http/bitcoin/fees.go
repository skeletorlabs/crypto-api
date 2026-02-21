package bitcoin

import (
	"crypto-api/internal/cache"
	api "crypto-api/internal/http"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/bitcoin"
	"encoding/json"
	"log"
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
			if err := json.NewEncoder(w).Encode(cached); err != nil {
				log.Printf("[http] failed to encode cached bitcoin fees: %v", err)
			}
			return
		}

		ctx := r.Context()
		fees, err := bitcoin.GetBitcoinFees(ctx)
		if err != nil {
			httpErr := api.MapError(err)
			api.JSONError(w, httpErr.Status, httpErr.Message)
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

		cache.Set(c, cacheKey, resp, cache.TTLBitcoinFees)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[http] failed to encode bitcoin fees response: %v", err)
		}
	}
}
