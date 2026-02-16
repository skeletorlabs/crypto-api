package market

import (
	"context"
	"crypto-api/internal/sources"
	"time"
)

func GetProtocols(ctx context.Context) ([]DefillamaProtocol, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var protocols []DefillamaProtocol
	url := "https://api.llama.fi/protocols"

	if err := sources.FetchJSON(ctx, url, &protocols); err != nil {
		return nil, err
	}

	return protocols, nil
}
