package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CoinGeckoPriceResponse map[string]map[string]float64

func GetPriceUSD(symbol string) (float64, error) {
	symbol = strings.ToLower(symbol)

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd",
		symbol,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("coingecko returned status %d", resp.StatusCode)
	}

	var data CoinGeckoPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, ok := data[symbol]["usd"]
	if !ok {
		return 0, fmt.Errorf("price not found for %s", symbol)
	}

	return price, nil
}
