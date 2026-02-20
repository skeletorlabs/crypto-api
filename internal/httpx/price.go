package httpx

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/engine/intelligence"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
	"crypto-api/internal/storage/repositories"
)

func IntelligencePriceHandler(c *cache.MemoryCache, repo *repositories.IntelligenceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/intelligence/price/"))
		if token == "" || token == "price" {
			JSONError(w, http.StatusBadRequest, "missing token")
			return
		}

		if !intelligence.IsSupportedAsset(token) {
			JSONError(w, http.StatusNotFound, "token not supported by intelligence engine")
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

		JSONError(w, http.StatusServiceUnavailable, "intelligence data not yet synchronized")
	}
}

func MarketPriceHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/market/price/"))
		if token == "" || token == "price" {
			JSONError(w, http.StatusBadRequest, "missing token")
			return
		}

		cacheKey := cache.KeyMarketPrice(token)
		updatedAt := time.Now().UTC()

		if price, ok := cache.Get[float64](c, cacheKey); ok {
			sendPriceResponse(w, token, price, true, updatedAt)
			return
		}

		ctx := r.Context()
		price, err := market.GetPriceUSD(ctx, token)
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		cache.Set(c, cacheKey, price, cache.TTLBitcoinPrice)
		sendPriceResponse(w, token, price, false, updatedAt)
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
