package intelligence

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/intelligence"
	api "crypto-api/internal/http"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func IntelligencePriceHandler(c *cache.MemoryCache, repo *repositories.IntelligenceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/intelligence/price/"))
		if token == "" || token == "price" {
			api.JSONError(w, http.StatusBadRequest, "missing token")
			return
		}

		if !intelligence.IsSupportedAsset(token) {
			api.JSONError(w, http.StatusNotFound, "token not supported by intelligence engine")
			return
		}

		cacheKey := cache.KeyIntelligencePrice(token)

		if data, ok := cache.Get[models.IntelligencePrice](c, cacheKey); ok {
			sendPriceResponse(w, token, data.Price, true, data.CreatedAt)
			return
		}

		snapshot, err := repo.GetLatest(r.Context())
		if err == nil && snapshot != nil {
			data := models.IntelligencePrice{
				Price:     snapshot.PriceUSD,
				CreatedAt: snapshot.CreatedAt,
			}

			cache.Set(c, cacheKey, data, cache.TTLBitcoinPrice)
			sendPriceResponse(w, token, data.Price, false, data.CreatedAt)
			return
		}

		api.JSONError(w, http.StatusServiceUnavailable, "intelligence data not yet synchronized")
	}
}

func sendPriceResponse(w http.ResponseWriter, token string, price float64, cached bool, updatedAt time.Time) {
	resp := models.PriceResponse{
		Meta: models.Meta{
			UpdatedAt: updatedAt,
			Cached:    cached,
		},
		Token: token,
		USD:   price,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[http] failed to encode price response for %s: %v", token, err)
	}
}
