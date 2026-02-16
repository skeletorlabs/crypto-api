package market

import (
	"context"
	"crypto-api/internal/sources"
	"fmt"
	"strings"
	"time"
)

func priceFromCoingecko(ctx context.Context, token string) (float64, error) {
	token = strings.ToLower(token)
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", token)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var data CoingeckoPriceResponse

	if err := sources.FetchJSON(ctx, url, &data); err != nil {
		return 0, err
	}

	price, ok := data[token]["usd"]
	if !ok {
		return 0, fmt.Errorf("price not found for %s", token)
	}

	return price, nil
}
