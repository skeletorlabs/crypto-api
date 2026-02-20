package market

import (
	"context"
	"crypto-api/internal/sources"
)

func GetProtocols(ctx context.Context) ([]DefillamaProtocol, error) {
	var protocols []DefillamaProtocol
	url := "https://api.llama.fi/protocols"

	if err := sources.FetchJSON(ctx, url, &protocols); err != nil {
		return nil, err
	}

	return protocols, nil
}
