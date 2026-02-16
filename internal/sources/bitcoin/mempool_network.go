package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NetworkRawData struct {
	BlockHeight  int64
	HashrateTHs  float64
	Difficulty   float64
	AvgBlockTime float64
}

type mempoolBlock struct {
	Timestamp int64 `json:"timestamp"`
}

type mempoolHashratePoint struct {
	Timestamp   int64   `json:"timestamp"`
	AvgHashrate float64 `json:"avgHashrate"`
}

type mempoolHashrate struct {
	Hashrates         []mempoolHashratePoint `json:"hashrates"`
	CurrentDifficulty float64                `json:"currentDifficulty"`
}

func GetBitcoinNetwork(ctx context.Context) (*NetworkRawData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	height, err := getBitcoinBlockHeight(ctx)
	if err != nil {
		return nil, err
	}

	hashrate, difficulty, err := getBitcoinHashrateTHs(ctx)
	if err != nil {
		return nil, err
	}

	avgTime, err := getBitcoinAvgBlockTime(ctx)
	if err != nil {
		return nil, err
	}

	return &NetworkRawData{
		BlockHeight:  height,
		HashrateTHs:  hashrate,
		Difficulty:   difficulty,
		AvgBlockTime: avgTime,
	}, nil
}

func getBitcoinBlockHeight(ctx context.Context) (int64, error) {
	url := "https://mempool.space/api/blocks/tip/height"

	body, err := sources.FetchRaw(ctx, url)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(strings.TrimSpace(string(body)), 10, 64)

}

func getBitcoinHashrateTHs(ctx context.Context) (float64, float64, error) {
	url := "https://mempool.space/api/v1/mining/hashrate/1y"

	var hrResp mempoolHashrate
	if err := sources.FetchJSON(ctx, url, &hrResp); err != nil {
		return 0, 0, err
	}

	if len(hrResp.Hashrates) == 0 {
		return 0, 0, fmt.Errorf("empty hashrate data")
	}

	last := hrResp.Hashrates[len(hrResp.Hashrates)-1].AvgHashrate
	return last / 1e12, hrResp.CurrentDifficulty, nil
}

func getBitcoinAvgBlockTime(ctx context.Context) (float64, error) {
	url := "https://mempool.space/api/blocks"

	var blocks []mempoolBlock
	if err := sources.FetchJSON(ctx, url, &blocks); err != nil {
		return 0, err
	}

	if len(blocks) < 2 {
		return 0, fmt.Errorf("insufficient blocks for average")
	}

	var sum float64
	count := len(blocks) - 1
	for i := 0; i < count; i++ {
		sum += float64(blocks[i].Timestamp - blocks[i+1].Timestamp)
	}

	return sum / float64(count), nil
}
