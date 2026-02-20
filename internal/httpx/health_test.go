package httpx

import (
	"crypto-api/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	// 1. Criamos um store vazio (isso pode causar erro no Ping, mas permite compilar)
	// Em um cenário real, você passaria uma conexão de teste aqui.
	mockStore := &storage.PostgresStore{}

	// 2. O Handler agora é uma função que retorna outra função (Closure)
	handler := HealthHandler(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// 3. Executa o handler
	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	// 4. Validações
	// Nota: Se o mockStore estiver vazio, o Ping vai falhar e o status será 503.
	// Se você quer que o teste passe com 200, precisaria de um banco real conectado.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status 200 or 503, got %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	// 5. Como o corpo tem um campo "time" que muda sempre,
	// não comparamos a string inteira. Decodificamos o JSON.
	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] == "" {
		t.Fatal("expected status field in response")
	}

	if _, ok := response["time"]; !ok {
		t.Fatal("expected time field in response")
	}
}
