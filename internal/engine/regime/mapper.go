package regime

import (
	"strings"

	"crypto-api/internal/models"
)

func ToPhaseResponse(asset string, r Result) models.PhaseResponse {
	return models.PhaseResponse{
		Asset:      asset,
		Phase:      mapRegime(r.Regime),
		Confidence: r.Confidence,
		Status:     models.Status(toTitle(string(r.Momentum))),
	}
}

func mapRegime(r Regime) models.Phase {
	switch r {
	case RegimeExpansion:
		return models.PhaseExpansion
	case RegimeContraction:
		return models.PhaseContraction
	default:
		return models.PhaseTransition
	}
}

func toTitle(s string) string {
	if s == "" {
		return s
	}

	s = strings.ToLower(s)
	return strings.ToUpper(string(s[0])) + s[1:]
}
