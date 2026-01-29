package httpx

import (
	"crypto-api/internal/cache"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBitcoinMempoolHandler_Cache(t *testing.T) {
	c := cache.NewMemoryCache()

	req := httptest.NewRequest(http.MethodGet, "/bitcoin/mempool", nil)
	rr := httptest.NewRecorder()

	handler := GetBitcoinMempoolHandler(c)
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
