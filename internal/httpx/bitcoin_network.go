package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func BitcoinNetworkHandler(c *cache.MemoryCache, repo *repositories.NetworkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyBitcoinNetwork

		// 1. Cache first
		if cached, ok := cache.Get[models.BitcoinNetworkResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			if err := json.NewEncoder(w).Encode(cached); err != nil {
				log.Printf("[http] failed to encode cached network response: %v", err)
			}
			return
		}

		// 2. Database fallback
		// 2. Database fallback
		data, err := repo.GetLatest(r.Context())
		if err != nil {
			JSONError(w, http.StatusServiceUnavailable, "Network data currently unavailable")
			return
		}

		now := time.Now().UTC()

		// Recalculate Halving
		data.Halving = halving.Compute(
			int(data.BlockHeight),
			data.AvgBlockTimeSeconds/60,
			now,
		)

		// Fetch previous record to compute trend
		prev, err := repo.GetPrevious(r.Context())

		hasPrev := err == nil && prev != nil
		prevAvg := 0.0
		if hasPrev {
			prevAvg = prev.AvgBlockTimeSeconds
		}

		// Recalculate Trend
		data.Trend = valuation.CalculateTrend(
			data.AvgBlockTimeSeconds,
			prevAvg,
			hasPrev,
		)

		data.Meta.Cached = false

		cache.Set(c, cacheKey, *data, cache.TTLNetworkStats)

		json.NewEncoder(w).Encode(data)
	}
}
