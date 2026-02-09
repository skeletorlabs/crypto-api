package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/trend"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/bitcoin"
	"encoding/json"
	"net/http"
	"time"
)

var bitcoinNetworkTrendBuffer = trend.NewBuffer(20)

func BitcoinNetworkHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyBitcoinNetwork

		if cached, ok := cache.Get[models.BitcoinNetworkResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			json.NewEncoder(w).Encode(cached)
			return
		}

		ctx := r.Context()
		data, err := bitcoin.GetBitcoinNetwork(ctx)

		if err != nil {
			httpErr := MapError((err))
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		snap := trend.Snapshot{
			Timestamp:       time.Now().UTC(),
			AvgBlockTimeSec: data.AvgBlockTime,
		}

		bitcoinNetworkTrendBuffer.Add(snap)
		trendStatus := trend.ComputeTrend(
			bitcoinNetworkTrendBuffer.All(),
		)
		halvingState := halving.Compute(int(data.BlockHeight), data.AvgBlockTime/60)

		resp := models.BitcoinNetworkResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			BlockHeight:         data.BlockHeight,
			HashrateTHs:         data.HashrateTHs,
			Difficulty:          data.Difficulty,
			AvgBlockTimeSeconds: data.AvgBlockTime,
			Trend:               trendStatus,
			Halving:             halvingState,
		}

		cache.Set(c, cacheKey, resp, 30*time.Second)

		json.NewEncoder(w).Encode(resp)
	}
}
