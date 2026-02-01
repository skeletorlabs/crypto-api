package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

func GetBitcoinMempoolHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if cached, ok := c.Get("bitcoin:mempool"); ok {
			resp := cached.(models.BitcoinMempoolResponse)
			resp.Cached = true
			json.NewEncoder(w).Encode(resp)
			return
		}

		stats, err := sources.GetBitcoinMempool()
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
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

		c.Set("bitcoin:mempool", resp, 30*time.Second)

		json.NewEncoder(w).Encode(resp)
	}
}
