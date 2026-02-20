package valuation

import (
	"math"

	"crypto-api/internal/models"
)

// CalculateNetworkHealth computes a 0–100 score based on deviation
// from the 10-minute Bitcoin block target (600 seconds).
func CalculateNetworkHealth(avgBlockTime float64) int {
	if avgBlockTime <= 0 {
		return 0
	}

	const target = 600.0 // seconds

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
	previous *models.IntelligenceSnapshot,
) models.Status {

	if previous == nil {
		return models.TrendStable
	}

	const epsilon = 30.0 // seconds tolerance

	diff := currentAvgTime - previous.AvgBlockTime

	if diff > epsilon {
		return models.TrendWorsening
	}
	if diff < -epsilon {
		return models.TrendImproving
	}

	return models.TrendStable
}
