package regime

type Regime string

const (
	RegimeExpansion   Regime = "EXPANSIONARY"
	RegimeContraction Regime = "CONTRACTIONARY"
	RegimeTransition  Regime = "TRANSITIONAL"
)

type Momentum string

const (
	MomentumImproving Momentum = "IMPROVING"
	MomentumWeakening Momentum = "WEAKENING"
	MomentumNeutral   Momentum = "NEUTRAL"
)

// Result represents the structural regime evaluation output.
type Result struct {
	Regime     Regime
	Score      float64
	Confidence float64
	Momentum   Momentum
}
