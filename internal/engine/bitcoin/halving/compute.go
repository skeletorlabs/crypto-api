package halving

import (
	"crypto-api/internal/config"
	"math"
	"time"
)

// Compute calculates halving cycle state based on
// current block height and average block time (in minutes).
func Compute(currentBlock int, avgBlockTime float64, baseTime time.Time) State {
	halvingInterval := config.HalvingInterval
	const BitcoinInitialSubsidy = 50.0 // Initial BTC block reward

	epoch := currentBlock / halvingInterval

	nextHalvingBlock := (epoch + 1) * halvingInterval
	blocksRemaining := nextHalvingBlock - currentBlock

	startBlock := epoch * halvingInterval

	progress := float64(currentBlock-startBlock) / float64(halvingInterval) * 100
	progressPercent := math.Round(progress*100) / 100

	remainingMinutes := float64(blocksRemaining) * avgBlockTime
	estimatedDate :=
		baseTime.UTC().Add(time.Duration(remainingMinutes) * time.Minute)

	currentSubsidy := BitcoinInitialSubsidy / math.Pow(2, float64(epoch))
	nextSubsidy := currentSubsidy / 2

	return State{
		CurrentBlock:     currentBlock,
		NextHalvingBlock: nextHalvingBlock,
		BlocksRemaining:  blocksRemaining,
		ProgressPercent:  progressPercent,
		EstimatedDate:    estimatedDate,
		CurrentSubsidy:   currentSubsidy,
		NextSubsidy:      nextSubsidy,
	}
}
