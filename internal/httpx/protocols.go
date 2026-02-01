package httpx

import (
	"encoding/json"
	"net/http"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/filters"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
)

func ProtocolsHandler(protocolsCache *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		chain := r.URL.Query().Get("chain")
		category := r.URL.Query().Get("category")

		var resp models.StandardResponse[[]models.ProtocolResponse]

		if cached, ok := protocolsCache.Get("all"); ok {
			resp = cached.(models.StandardResponse[[]models.ProtocolResponse])
			resp.Meta.Cached = true
		} else {
			data, err := sources.GetProtocols()
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

			protocolsCache.Set("all", resp, 5*time.Minute)
		}

		resp.Data = filters.FilterProtocols(resp.Data, chain, category)
		json.NewEncoder(w).Encode(resp)
	}
}
