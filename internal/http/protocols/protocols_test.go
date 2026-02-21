package protocols

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-api/internal/cache"
)

func TestProtocolsHandler_OK(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := ProtocolsHandler(c)

	req := httptest.NewRequest(
		http.MethodGet,
		"/protocols?chain=Ethereum",
		nil,
	)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if _, ok := body["meta"]; !ok {
		t.Fatalf("expected 'meta' field in response")
	}

	if _, ok := body["data"]; !ok {
		t.Fatalf("expected 'data' field in response")
	}
}
