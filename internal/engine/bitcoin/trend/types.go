package trend

type Status string

const (
	TrendImproving Status = "Improving"
	TrendStable    Status = "Stable"
	TrendWorsening Status = "Worsening"
)
