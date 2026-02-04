package trend

func ComputeTrend(snaps []Snapshot) Status {
	if len(snaps) < 2 {
		return TrendStable
	}

	prev := snaps[len(snaps)-2]
	curr := snaps[len(snaps)-1]

	diff := curr.AvgBlockTimeSec - prev.AvgBlockTimeSec

	const epsilon = 30.0 // seconds

	switch {
	case diff > epsilon:
		return TrendWorsening
	case diff < -epsilon:
		return TrendImproving
	default:
		return TrendStable
	}
}
