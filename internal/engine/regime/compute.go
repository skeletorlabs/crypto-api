package regime

import "crypto-api/internal/models"

// Compute evaluates the structural regime based on an IntelligenceSnapshot.
// Placeholder logic until the real scoring model is implemented.
func Compute(snapshot models.IntelligenceSnapshot) Result {
	return Result{
		Regime:     RegimeExpansion,
		Score:      0.74,
		Confidence: 0.74,
		Momentum:   MomentumImproving,
	}
}
