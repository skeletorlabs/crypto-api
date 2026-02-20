package httpx

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/filters"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
)

func ProtocolsHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		chain := r.URL.Query().Get("chain")
		category := r.URL.Query().Get("category")
		cacheKey := cache.KeyMarketProtocols

		var resp models.StandardResponse[[]models.ProtocolResponse]

		if cached, ok := cache.Get[models.StandardResponse[[]models.ProtocolResponse]](c, cacheKey); ok {
			resp = cached
			resp.Meta.Cached = true
		} else {
			ctx := r.Context()
			data, err := market.GetProtocols(ctx)
			if err != nil {
				httpErr := MapError(err)
				JSONError(w, httpErr.Status, httpErr.Message)
				return
			}

			protocols := make([]models.ProtocolResponse, 0, len(data))
			for _, protocol := range data {
				protocols = append(protocols, models.ProtocolResponse{
					Name:     protocol.Name,
					Slug:     protocol.Slug,
					TVL:      protocol.TVL,
					Chain:    protocol.Chain,
					Category: protocol.Category,
				})
			}

			resp = models.StandardResponse[[]models.ProtocolResponse]{
				Meta: models.Meta{
					UpdatedAt: time.Now().UTC(),
					Cached:    false,
				},
				Data: protocols,
			}

			cache.Set(c, cacheKey, resp, cache.TTLNetworkStats)
		}

		resp.Data = filters.FilterProtocols(resp.Data, chain, category)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[http] failed to encode protocols response: %v", err)
		}
	}
}
