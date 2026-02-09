package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/macro"
	"encoding/json"
	"net/http"
	"time"
)

// This handler is responsible for fetching macroeconomic data (like M2 supply) and caching it.
// For now is only M2, but we can easily expand it to include more indicators in the future.
func MacroHandler(c *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		cacheKey := cache.KeyMacroM2Supply

		if cached, ok := cache.Get[models.MacroResponse](c, cacheKey); ok {
			cached.Meta.Cached = true
			json.NewEncoder(w).Encode(cached)
			return
		}

		ctx := r.Context()
		m2Value, date, err := macro.GetM2Supply(ctx)
		if err != nil {
			httpErr := MapError(err)
			JSONError(w, httpErr.Status, httpErr.Message)
			return
		}

		resp := models.MacroResponse{
			Meta: models.Meta{
				UpdatedAt: time.Now().UTC(),
				Cached:    false,
			},
			M2Supply: models.M2Details{
				Value:    m2Value,
				Unit:     "Billions of Dollars",
				DateTime: date,
			},
		}

		cache.Set(c, cacheKey, resp, 24*time.Hour)
		json.NewEncoder(w).Encode(resp)
	}
}
