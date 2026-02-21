package market

import (
	"crypto-api/internal/cache"
	api "crypto-api/internal/http"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func MarketPriceHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/market/price/"))
		if token == "" || token == "price" {
			api.JSONError(w, http.StatusBadRequest, "missing token")
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
			httpErr := api.MapError(err)
			api.JSONError(w, httpErr.Status, httpErr.Message)
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
