package correlation

import (
	"errors"
	"math"
)

// ComputePearsonCorrelation calculates the Pearson correlation coefficient
// between two datasets (e.g., BTC prices and M2 Supply history).
func pearsonCorrelation(datasetA []float64, datasetB []float64) (float64, error) {
	n := len(datasetA)
	if n != len(datasetB) {
		return 0, errors.New("datasets must have the same length")
	}
	if n == 0 {
		return 0, errors.New("datasets cannot be empty")
	}

	var sumA, sumB, sumAB, sumA2, sumB2 float64

	for i := 0; i < n; i++ {
		sumA += datasetA[i]
		sumB += datasetB[i]
		sumAB += datasetA[i] * datasetB[i]
		sumA2 += datasetA[i] * datasetA[i]
		sumB2 += datasetB[i] * datasetB[i]
	}

	numerator := (float64(n) * sumAB) - (sumA * sumB)
	denominator := math.Sqrt((float64(n)*sumA2 - (sumA * sumA)) * (float64(n)*sumB2 - (sumB * sumB)))

	if denominator == 0 {
		return 0, nil
	}

	return numerator / denominator, nil
}
