package market

import (
	"crypto-api/internal/cache"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMarketPriceHandler_MissingToken(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := MarketPriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/market/price/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestMarketPriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	const fakePrice = 123.45
	token := "bitcoin"
	cacheKey := cache.KeyMarketPrice(token)

	cache.Set(c, cacheKey, fakePrice, cache.TTLBitcoinPrice)

	handler := MarketPriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/market/price/bitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"cached":true`) {
		t.Errorf("expected cached=true, got %s", body)
	}
	if !strings.Contains(body, "123.45") {
		t.Errorf("expected price 123.45, got %s", body)
	}
}
