package filters

import "crypto-api/internal/models"

func FilterProtocols(protocols []models.ProtocolResponse, chain string, category string) []models.ProtocolResponse {
	if chain == "" && category == "" {
		return protocols
	}

	filtered := make([]models.ProtocolResponse, 0)
	for _, p := range protocols {
		if chain != "" && p.Chain != chain {
			continue
		}
		if category != "" && p.Category != category {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
}
