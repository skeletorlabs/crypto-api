package httpx

import (
	"encoding/json"
	"net/http"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
)

func ChainsHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyMarketChains

		if cached, ok := cache.Get[models.StandardResponse[[]models.ChainResponse]](c, cacheKey); ok {
			cached.Meta.Cached = true
			json.NewEncoder(w).Encode(cached)
			return
		}

		ctx := r.Context()
		chains, err := market.GetChains(ctx)
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		response := make([]models.ChainResponse, 0, len(chains))
		for _, chain := range chains {
			response = append(response, models.ChainResponse{
				Name:   chain.Name,
				TVL:    chain.TVL,
				Symbol: chain.TokenSymbol,
			})
		}

		resp := models.StandardResponse[[]models.ChainResponse]{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			Data: response,
		}

		cache.Set(c, cacheKey, resp, 5*time.Minute)
		json.NewEncoder(w).Encode(resp)
	}
}
