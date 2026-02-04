package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/trend"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

var bitcoinNetworkTrendBuffer = trend.NewBuffer(20)

func BitcoinNetworkHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if cached, ok := c.Get("network"); ok {
			resp := cached.(models.BitcoinNetworkResponse)
			resp.Cached = true
			json.NewEncoder(w).Encode(resp)
			return
		}

		height, hashrate, difficulty, avgtime, err := sources.GetBitcoinNetwork()

		if err != nil {
			httpErr := MapError((err))
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		snap := trend.Snapshot{
			Timestamp:       time.Now().UTC(),
			AvgBlockTimeSec: avgtime,
		}

		bitcoinNetworkTrendBuffer.Add(snap)
		trendStatus := trend.ComputeTrend(
			bitcoinNetworkTrendBuffer.All(),
		)
		halvingState := halving.Compute(int(height), avgtime/60)

		resp := models.BitcoinNetworkResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			BlockHeight:         height,
			HashrateTHs:         hashrate,
			Difficulty:          difficulty,
			AvgBlockTimeSeconds: avgtime,
			Trend:               trendStatus,
			Halving:             halvingState,
		}

		c.Set("network", resp, 30*time.Second)

		json.NewEncoder(w).Encode(resp)
	}
}
