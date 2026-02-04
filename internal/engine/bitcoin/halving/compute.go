package halving

import (
	"math"
	"time"
)

func Compute(currentBlock int, avgBlockTime float64) State {
	const halvingInterval = 210000

	epoch := currentBlock / halvingInterval
	nextHalvingBlock := (epoch + 1) * halvingInterval
	blocksRemaining := nextHalvingBlock - currentBlock

	startBlock := epoch * halvingInterval
	progress := float64(currentBlock-startBlock) / float64(halvingInterval) * 100
	progressPercent := math.Round(progress*100) / 100

	remainingMinutes := float64(blocksRemaining) * avgBlockTime
	estimatedDate :=
		time.Now().UTC().Add(time.Duration(remainingMinutes) * time.Minute)

	return State{
		CurrentBlock:     currentBlock,
		NextHalvingBlock: nextHalvingBlock,
		BlocksRemaining:  blocksRemaining,
		ProgressPercent:  progressPercent,
		EstimatedDate:    estimatedDate,
	}
}
