package chains

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crypto-api/internal/cache"
	api "crypto-api/internal/http"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
)

func ChainsHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyMarketChains

		if cached, ok := cache.Get[models.StandardResponse[[]models.ChainResponse]](c, cacheKey); ok {
			cached.Meta.Cached = true
			if err := json.NewEncoder(w).Encode(cached); err != nil {
				log.Printf("[http] failed to encode cached chains response: %v", err)
			}
			return
		}

		ctx := r.Context()
		chains, err := market.GetChains(ctx)
		if err != nil {
			httpErr := api.MapError(err)
			api.JSONError(w, httpErr.Status, httpErr.Message)
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

		cache.Set(c, cacheKey, resp, cache.TTLNetworkStats)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[http] failed to encode chains response: %v", err)
		}
	}
}
