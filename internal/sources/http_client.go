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

var HttpClient HTTPDoer = &http.Client{
	Timeout: 45 * time.Second,
}

func FetchJSON[T any](ctx context.Context, url string, target *T) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpstreamTimeout, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upstream error: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func FetchRaw(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
