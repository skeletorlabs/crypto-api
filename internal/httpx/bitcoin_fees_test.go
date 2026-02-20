package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBitcoinFeesHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	cachedKey := cache.KeyBitcoinFees

	cache.Set(c, cachedKey, models.BitcoinFeesResponse{
		Low:    10,
		Medium: 20,
		High:   30,
	}, cache.TTLBitcoinPrice)

	handler := BitcoinFeesHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/bitcoin/fees", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", rr.Body.String())
	}
}
