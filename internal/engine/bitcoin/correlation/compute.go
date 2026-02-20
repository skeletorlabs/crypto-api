package correlation

import (
	"fmt"
	"sort"
)

// Compute aligns M2 and BTC data by date and calculates their Pearson correlation.
// Input slices are not mutated.
func Compute(m2Data []DataPoint, btcData []DataPoint) (*Result, error) {
	if len(m2Data) == 0 || len(btcData) == 0 {
		return nil, fmt.Errorf("cannot compute correlation with empty datasets")
	}

	// Defensive copy to avoid mutating caller slices
	m2 := append([]DataPoint(nil), m2Data...)
	btc := append([]DataPoint(nil), btcData...)

	// 1. Sort data (Oldest -> Newest)
	sort.Slice(m2, func(i, j int) bool { return m2[i].Date.Before(m2[j].Date) })
	sort.Slice(btc, func(i, j int) bool { return btc[i].Date.Before(btc[j].Date) })

	// 2. Map M2 for O(1) lookup
	m2Map := make(map[string]float64, len(m2))
	for _, p := range m2 {
		m2Map[p.Date.Format("2006-01-02")] = p.Value
	}

	var alignedM2 []float64
	var alignedBTC []float64
	var lastKnownM2 float64

	if len(m2) > 0 {
		lastKnownM2 = m2[0].Value
	}

	// 3. Alignment with Forward Fill logic
	for _, btcPoint := range btc {
		dateStr := btcPoint.Date.Format("2006-01-02")

		if val, ok := m2Map[dateStr]; ok {
			lastKnownM2 = val
		}

		if lastKnownM2 > 0 {
			alignedM2 = append(alignedM2, lastKnownM2)
			alignedBTC = append(alignedBTC, btcPoint.Value)
		}
	}

	if len(alignedM2) < 2 {
		return nil, fmt.Errorf(
			"insufficient aligned data points (found %d, need at least 2)",
			len(alignedM2),
		)
	}

	// 4. Statistical Calculation
	coef, err := pearsonCorrelation(alignedM2, alignedBTC)
	if err != nil {
		return nil, fmt.Errorf("pearson math failed: %w", err)
	}

	return &Result{
		Coefficient: coef,
		SampleCount: len(alignedM2),
		StartDate:   btc[0].Date,
		EndDate:     btc[len(btc)-1].Date,
	}, nil
}
