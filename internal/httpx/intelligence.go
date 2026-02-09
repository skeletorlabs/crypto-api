package httpx

import (
	"context"
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/macro"
	"crypto-api/internal/sources/market"
	"encoding/json"
	"net/http"
	"time"
)

func ValuationHandler(marketCache, macroCache, intelligenceCache *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1. Try to get pre-computed analysis
		if state, ok := cache.Get[valuation.State](intelligenceCache, cache.KeyIntelligenceValuation); ok {
			state.Meta.Cached = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(state)
			return
		}

		// 2. Search for BTC price
		btcPrice, ok := cache.Get[float64](marketCache, cache.KeyMarketPrice("bitcoin"))
		if !ok {
			p, err := market.GetPriceUSD(ctx, "bitcoin")
			if err != nil {
				http.Error(w, "Failed to fetch BTC price", http.StatusInternalServerError)
				return
			}
			btcPrice = p
			cache.Set(marketCache, cache.KeyMarketPrice("bitcoin"), btcPrice, 10*time.Minute)
		}

		// 3. Search for M2 Supply
		m2Value, ok := cache.Get[float64](macroCache, cache.KeyMacroM2Supply)
		if !ok {
			val, _, err := macro.GetM2Supply(ctx)
			if err != nil {
				http.Error(w, "Failed to fetch M2 supply", http.StatusInternalServerError)
				return
			}
			m2Value = val
			cache.Set(macroCache, cache.KeyMacroM2Supply, m2Value, 24*time.Hour)
		}

		// 4. Compute & save
		state := valuation.Compute(btcPrice, m2Value)
		state.Meta = models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
		}
		cache.Set(intelligenceCache, cache.KeyIntelligenceValuation, state, 10*time.Minute)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	}
}

func CorrelationHandler(intelligenceCache, marketCache, macroCache *cache.MemoryCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		// 1. Try Cache
		if result, ok := cache.Get[*correlation.Result](intelligenceCache, cache.KeyIntelligenceCorrelation); ok {
			result.Meta.Cached = true
			json.NewEncoder(w).Encode(result)
			return
		}

		m2History, err := getM2HistoryData(ctx, macroCache)
		if err != nil {
			http.Error(w, "M2 data unavailable", http.StatusBadGateway)
			return
		}

		btcHistory, err := getBTCHistoryData(ctx, marketCache)
		if err != nil {
			http.Error(w, "BTC data unavailable", http.StatusBadGateway)
			return
		}

		result, err := correlation.Compute(m2History, btcHistory)
		if err != nil {
			http.Error(w, "Analysis computation failed", http.StatusInternalServerError)
			return
		}

		result.Meta = models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
		}

		cache.Set(intelligenceCache, cache.KeyIntelligenceCorrelation, result, 24*time.Hour)
		json.NewEncoder(w).Encode(result)
	}
}

func getM2HistoryData(ctx context.Context, macroCache *cache.MemoryCache) ([]correlation.DataPoint, error) {
	key := cache.KeyMacroM2History
	if data, ok := cache.Get[[]correlation.DataPoint](macroCache, key); ok {
		return data, nil
	}
	h, err := macro.GetM2History(ctx, 100)
	if err != nil {
		return nil, err
	}
	cache.Set(macroCache, key, h, 24*time.Hour)
	return h, nil
}

func getBTCHistoryData(ctx context.Context, marketCache *cache.MemoryCache) ([]correlation.DataPoint, error) {
	key := cache.KeyMarketHistory("bitcoin")
	if data, ok := cache.Get[[]correlation.DataPoint](marketCache, key); ok {
		return data, nil
	}
	h, err := market.GetBTCPriceHistory(ctx, 730)
	if err != nil {
		return nil, err
	}
	cache.Set(marketCache, key, h, 12*time.Hour)
	return h, nil
}
