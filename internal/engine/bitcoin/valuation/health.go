package valuation

import (
	"math"

	"crypto-api/internal/config"
	"crypto-api/internal/models"
)

// CalculateNetworkHealth computes a 0–100 score based on deviation
// from the 10-minute Bitcoin block target (600 seconds).
func CalculateNetworkHealth(avgBlockTime float64) int {
	if avgBlockTime <= 0 {
		return 0
	}

	target := config.BitcoinTargetBlockSeconds
	diff := math.Abs(avgBlockTime - target)
	deviationPercent := (diff / target) * 100
	score := 100.0 - (deviationPercent * 2)
	if score > 100 {
		return 100
	}
	if score < 0 {
		return 0
	}

	return int(score)
}

// CalculateTrend determines block time direction relative to
// the previous intelligence snapshot using a fixed epsilon threshold.
func CalculateTrend(
	currentAvgTime float64,
	previousAvgTime float64,
	hasPrevious bool,
) models.Status {

	if !hasPrevious {
		return models.TrendStable
	}

	epsilon := config.TrendEpsilonSeconds
	diff := currentAvgTime - previousAvgTime

	if diff > epsilon {
		return models.TrendWorsening
	}
	if diff < -epsilon {
		return models.TrendImproving
	}

	return models.TrendStable
}
