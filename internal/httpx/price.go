package httpx

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
)

func PriceHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/price/")
		if token == "" {
			JSONError(w, http.StatusBadRequest, "missing token")
			return
		}
		token = strings.ToLower(token)
		cacheKey := cache.KeyMarketPrice(token)

		if cached, ok := cache.Get[models.PriceResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			json.NewEncoder(w).Encode(cached)
			return
		}

		ctx := r.Context()
		price, err := market.GetPriceUSD(ctx, token)
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		resp := models.PriceResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			Token: token,
			USD:   price,
		}

		cache.Set(c, cacheKey, resp, 30*time.Second)
		json.NewEncoder(w).Encode(resp)
	}
}
