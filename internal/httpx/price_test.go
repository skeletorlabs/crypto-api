package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"crypto-api/internal/cache"
)

func TestPriceHandler_MissingSymbol(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := PriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/price/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	const fakePrice = 123.45
	c.Set("bitcoin", fakePrice, time.Minute)

	handler := PriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/price/bitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", rr.Body.String())
	}
}
