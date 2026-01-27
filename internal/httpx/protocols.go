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

		var protocols []models.ProtocolResponse

		if cached, ok := protocolsCache.Get("all"); ok {
			protocols = cached.([]models.ProtocolResponse)
		} else {
			data, err := sources.GetProtocols()
			if err != nil {
				JSONError(w, http.StatusInternalServerError, "failed to fetch protocols")
				return
			}

			protocols = make([]models.ProtocolResponse, 0, len(data))
			for _, protocol := range data {
				protocols = append(protocols, models.ProtocolResponse{
					Name:     protocol.Name,
					Slug:     protocol.Slug,
					TVL:      protocol.TVL,
					Chain:    protocol.Chain,
					Category: protocol.Category,
				})
			}

			protocolsCache.Set("all", protocols, 5*time.Minute)
		}

		filtered := filters.FilterProtocols(protocols, chain, category)
		json.NewEncoder(w).Encode(filtered)
	}
}
