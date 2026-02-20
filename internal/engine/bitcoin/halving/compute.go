package halving

import (
	"math"
	"time"
)

// Compute calculates halving cycle state based on
// current block height and average block time (in minutes).
func Compute(currentBlock int, avgBlockTime float64) State {
	const halvingInterval = 210000
	const initialSubsidy = 50.0 // Initial BTC block reward

	epoch := currentBlock / halvingInterval

	nextHalvingBlock := (epoch + 1) * halvingInterval
	blocksRemaining := nextHalvingBlock - currentBlock

	startBlock := epoch * halvingInterval

	progress := float64(currentBlock-startBlock) / float64(halvingInterval) * 100
	progressPercent := math.Round(progress*100) / 100

	remainingMinutes := float64(blocksRemaining) * avgBlockTime
	estimatedDate :=
		time.Now().UTC().Add(time.Duration(remainingMinutes) * time.Minute)

	currentSubsidy := initialSubsidy / math.Pow(2, float64(epoch))
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
