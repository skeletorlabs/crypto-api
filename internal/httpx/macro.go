package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"log"
	"net/http"
)

func MacroHandler(c *cache.MemoryCache, repo *repositories.MacroRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cacheKey := cache.KeyMacroM2Supply

		// 1. Cache first
		if cached, ok := cache.Get[models.MacroResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			if err := json.NewEncoder(w).Encode(cached); err != nil {
				log.Printf("[http] failed to encode cached macro response: %v", err)
			}
			return
		}

		// 2. Database fallback
		supply, lastUpdate, err := repo.GetLatestM2(r.Context())
		if err != nil {
			JSONError(w, http.StatusServiceUnavailable, "Macro data currently unavailable")
			return
		}

		resp := models.MacroResponse{
			Meta: models.Meta{
				UpdatedAt: lastUpdate,
				Cached:    false,
			},
			M2Supply: models.M2Details{
				Value:    supply,
				Unit:     "Billions",
				DateTime: lastUpdate,
			},
		}

		cache.Set(c, cacheKey, resp, cache.TTLMacroData)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("[http] failed to encode macro response: %v", err)
		}
	}
}
