package correlation

// AlignAndCompute is the high-level function that uses compute.go logic
func Compute(m2Data []DataPoint, btcData []DataPoint) (*Result, error) {
	var alignedM2 []float64
	var alignedBTC []float64

	// Map for quick BTC lookup
	btcMap := make(map[string]float64)
	for _, p := range btcData {
		btcMap[p.Date.Format("2006-01-02")] = p.Value
	}

	// Alignment logic
	for _, m2 := range m2Data {
		dateStr := m2.Date.Format("2006-01-02")
		if btcVal, ok := btcMap[dateStr]; ok {
			alignedM2 = append(alignedM2, m2.Value)
			alignedBTC = append(alignedBTC, btcVal)
		}
	}

	coef, err := pearsonCorrelation(alignedM2, alignedBTC)
	if err != nil {
		return nil, err
	}

	return &Result{
		Coefficient: coef,
		SampleCount: len(alignedM2),
		StartDate:   m2Data[len(m2Data)-1].Date,
		EndDate:     m2Data[0].Date,
	}, nil
}
