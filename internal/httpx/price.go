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
		symbol := strings.TrimPrefix(r.URL.Path, "/price/")
		if symbol == "" {
			JSONError(w, http.StatusBadRequest, "missing symbol")
			return
		}
		symbol = strings.ToLower(symbol)

		if cached, ok := priceCache.Get(symbol); ok {
			resp := cached.(models.PriceResponse)
			resp.Meta.Cached = true
			json.NewEncoder(w).Encode(resp)
			return
		}

		price, err := sources.GetPriceUSD(symbol)
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
			Symbol: symbol,
			USD:    price,
		}

		priceCache.Set(symbol, resp, 30*time.Second)
		json.NewEncoder(w).Encode(resp)
	}
}
