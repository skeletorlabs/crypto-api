package httpx

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources"
)

func PriceHandler(priceCache *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.Path, "/price/")
		if token == "" {
			JSONError(w, http.StatusBadRequest, "missing token")
			return
		}
		token = strings.ToLower(token)

		if cached, ok := priceCache.Get(token); ok {
			resp := cached.(models.PriceResponse)
			resp.Meta.Cached = true
			json.NewEncoder(w).Encode(resp)
			return
		}

		price, err := sources.GetPriceUSD(token)
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

		priceCache.Set(token, resp, 30*time.Second)
		json.NewEncoder(w).Encode(resp)
	}
}
