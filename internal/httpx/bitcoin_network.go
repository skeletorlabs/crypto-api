package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

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

		resp := models.BitcoinNetworkResponse{
			BlockHeight:         height,
			HashrateTHs:         hashrate,
			Difficulty:          difficulty,
			AvgBlockTimeSeconds: avgtime,
		}

		c.Set("network", resp, 30*time.Second)

		json.NewEncoder(w).Encode(resp)
	}
}
