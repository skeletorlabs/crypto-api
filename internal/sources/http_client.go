package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Mantendo os 45s conforme sua necessidade para o FRED
var HttpClient HTTPDoer = &http.Client{
	Timeout: 45 * time.Second,
}

func FetchJSON[T any](ctx context.Context, url string, target *T) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Skeletor-Crypto-API/1.0")

	resp, err := HttpClient.Do(req)
	if err != nil {
		// Usando o erro que você já tem no errors.go
		return fmt.Errorf("%w: %v", ErrUpstreamTimeout, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Usando o erro que você já tem no errors.go
		return fmt.Errorf("%w: status %d", ErrUpstreamBadStatus, resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func FetchRaw(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Skeletor-Crypto-API/1.0")

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrUpstreamBadStatus, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
