package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValuationHandler_Success(t *testing.T) {
	// Setup caches
	marketCache := cache.NewMemoryCache()
	macroCache := cache.NewMemoryCache()
	intelCache := cache.NewMemoryCache()

	// Simulate BTC price and M2 supply in their respective caches
	cache.Set(marketCache, cache.KeyMarketPrice("bitcoin"), 65000.0, time.Minute)
	cache.Set(macroCache, cache.KeyMacroM2Supply, 21000.0, time.Minute)

	handler := ValuationHandler(marketCache, macroCache, intelCache)
	req := httptest.NewRequest(http.MethodGet, "/bitcoin/valuation", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp valuation.State
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Validations
	if resp.Meta.Cached {
		t.Error("expected cached=false on first hit")
	}
	if resp.BtcPrice != 65000.0 {
		t.Errorf("expected price 65000, got %f", resp.BtcPrice)
	}

	// Test if the result was stored in intelligenceCache
	if _, ok := cache.Get[valuation.State](intelCache, cache.KeyIntelligenceValuation); !ok {
		t.Error("result should be stored in intelligenceCache")
	}
}

func TestCorrelationHandler_FromCache(t *testing.T) {
	intelCache := cache.NewMemoryCache()
	marketCache := cache.NewMemoryCache()
	macroCache := cache.NewMemoryCache()

	// Simulate a pre-computed correlation result stored in the intelligence cache
	cachedResult := &correlation.Result{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
		},
		Coefficient: 0.85,
		SampleCount: 100,
	}

	cache.Set(intelCache, cache.KeyIntelligenceCorrelation, cachedResult, time.Minute)

	handler := CorrelationHandler(intelCache, marketCache, macroCache)
	req := httptest.NewRequest(http.MethodGet, "/bitcoin/correlation", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()

	if !strings.Contains(body, `"cached":true`) {
		t.Errorf("expected cached:true in JSON, got %s", body)
	}

	if !strings.Contains(body, `"coefficient":0.85`) {
		t.Errorf("expected coefficient:0.85 in JSON, got %s", body)
	}
}

func TestCorrelationHandler_Integration(t *testing.T) {
	// Guarantees that the flow of fetching historical data works and computes correlation correctly
	intelCache := cache.NewMemoryCache()
	marketCache := cache.NewMemoryCache()
	macroCache := cache.NewMemoryCache()

	// Mock of historical data in market and macro caches
	history := []correlation.DataPoint{
		{Date: time.Now(), Value: 100},
		{Date: time.Now().AddDate(0, 0, -1), Value: 110},
	}
	cache.Set(macroCache, cache.KeyMacroM2History, history, time.Minute)
	cache.Set(marketCache, cache.KeyMarketHistory("bitcoin"), history, time.Minute)

	handler := CorrelationHandler(intelCache, marketCache, macroCache)
	req := httptest.NewRequest(http.MethodGet, "/bitcoin/correlation", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
