package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

func GetBitcoinNetwork() (
	blockHeight int64,
	hashrateTHs float64,
	difficulty float64,
	avgBlockTime float64,
	err error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blockHeight, err = getBitcoinBlockHeight(ctx)
	if err != nil {
		return
	}

	hashrateTHs, difficulty, err = getBitcoinHashrateTHs(ctx)
	if err != nil {
		return
	}

	avgBlockTime, err = getBitcoinAvgBlockTime(ctx)
	if err != nil {
		return
	}

	return
}

func getBitcoinBlockHeight(ctx context.Context) (int64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://mempool.space/api/blocks/tip/height",
		nil,
	)
	if err != nil {
		return 0, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		return 0, sources.ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, sources.ErrUpstreamBadStatus
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	height, err := strconv.ParseInt(
		strings.TrimSpace(string(body)),
		10,
		64,
	)
	if err != nil {
		return 0, err
	}

	return height, nil
}

func getBitcoinHashrateTHs(ctx context.Context) (hashrateTHs float64, difficulty float64, err error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://mempool.space/api/v1/mining/hashrate/1y",
		nil,
	)
	if err != nil {
		return
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		err = sources.ErrUpstreamTimeout
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = sources.ErrUpstreamBadStatus
		return
	}

	var hrResp mempoolHashrate
	if err = json.NewDecoder(resp.Body).Decode(&hrResp); err != nil {
		return
	}

	if len(hrResp.Hashrates) == 0 {
		err = sources.ErrUpstreamBadStatus
		return
	}

	last := hrResp.Hashrates[len(hrResp.Hashrates)-1].AvgHashrate
	// last is returned in H/s by mempool.space => convert to TH/s
	hashrateTHs = last / 1e12
	difficulty = hrResp.CurrentDifficulty

	return
}

func getBitcoinAvgBlockTime(ctx context.Context) (float64, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://mempool.space/api/blocks",
		nil,
	)
	if err != nil {
		return 0, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		return 0, sources.ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, sources.ErrUpstreamBadStatus
	}

	var blocks []mempoolBlock
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		return 0, err
	}

	if len(blocks) < 2 {
		return 0, sources.ErrUpstreamBadStatus
	}

	var sum float64
	count := len(blocks) - 1
	for i := 0; i < count; i++ {
		sum += float64(blocks[i].Timestamp - blocks[i+1].Timestamp)
	}

	return sum / float64(count), nil
}
