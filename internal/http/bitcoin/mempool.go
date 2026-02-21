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

func BitcoinMempoolHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyBitcoinMempool

		if cached, ok := cache.Get[models.BitcoinMempoolResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			if err := json.NewEncoder(w).Encode(cached); err != nil {
				log.Printf("[http] failed to encode cached mempool response: %v", err)
			}
			return
		}

		ctx := r.Context()
		stats, err := bitcoin.GetBitcoinMempool(ctx)
		if err != nil {
			httpErr := api.MapError(err)
			api.JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		resp := models.BitcoinMempoolResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			Count:    stats.Count,
			VSize:    stats.VSize,
			TotalFee: stats.TotalFee,
		}

		cache.Set(c, cacheKey, resp, cache.TTLBitcoinFees)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[http] failed to encode mempool response: %v", err)
		}
	}
}
