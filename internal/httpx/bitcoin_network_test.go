package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBitcoinNetworkHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()

	cachedResp := models.BitcoinNetworkResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
		},
		BlockHeight:         100,
		HashrateTHs:         123.4,
		Difficulty:          999,
		AvgBlockTimeSeconds: 600,
	}

	cachedKey := cache.KeyBitcoinNetwork
	cache.Set(c, cachedKey, cachedResp, time.Minute)

	req := httptest.NewRequest(http.MethodGet, "/bitcoin/network", nil)
	rr := httptest.NewRecorder()

	handler := BitcoinNetworkHandler(c)
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", rr.Body.String())
	}
}
