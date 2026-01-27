package httpx

import (
	"encoding/json"
	"net/http"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
)

func ChainsHandler(chainsCache *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if cached, ok := chainsCache.Get("all"); ok {
			json.NewEncoder(w).Encode(cached)
			return
		}

		chains, err := sources.GetChains()
		if err != nil {
			JSONError(w, http.StatusInternalServerError, "failed to fetch chains")
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

		chainsCache.Set("all", response, 5*time.Minute)
		json.NewEncoder(w).Encode(response)
	}
}
