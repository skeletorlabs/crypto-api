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
			resp := models.PriceResponse{
				Symbol: symbol,
				USD:    cached.(float64),
				Cached: true,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		price, err := sources.GetPriceUSD(symbol)
		if err != nil {
			JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		priceCache.Set(symbol, price, 30*time.Second)

		resp := models.PriceResponse{
			Symbol: symbol,
			USD:    price,
			Cached: false,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
