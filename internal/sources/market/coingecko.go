package market

import (
	"context"
	"crypto-api/internal/sources"
	"errors"
	"fmt"
	"strings"
)

func priceFromCoingecko(ctx context.Context, token string) (float64, error) {
	token = strings.ToLower(token)
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", token)

	var data CoingeckoPriceResponse
	if err := sources.FetchJSON(ctx, url, &data); err != nil {
		// Se o erro for de status, verificamos se é rate limit
		if errors.Is(err, sources.ErrUpstreamBadStatus) {
			return 0, fmt.Errorf("coingecko rate limit or server error: %w", err)
		}
		return 0, err
	}

	priceMap, ok := data[token]
	if !ok {
		return 0, fmt.Errorf("token %s not found in coingecko response", token)
	}

	price, ok := priceMap["usd"]
	if !ok {
		return 0, fmt.Errorf("usd price not found for %s", token)
	}

	return price, nil
}
