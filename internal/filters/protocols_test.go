package filters

import (
	"testing"

	"crypto-api/internal/models"
)

func TestFilterProtocols_NoFilters(t *testing.T) {
	input := []models.ProtocolResponse{
		{Name: "Aave", Chain: "Ethereum", Category: "Lending"},
		{Name: "Uniswap", Chain: "Binance Smart Chain", Category: "DEX"},
	}

	result := FilterProtocols(input, "", "")
	if len(result) != 2 {
		t.Errorf("Expected 2 protocols, got %d", len(result))
	}
}

func TestFilterProtocols_ByChain(t *testing.T) {
	input := []models.ProtocolResponse{
		{Name: "Aave", Chain: "Ethereum", Category: "Lending"},
		{Name: "PancakeSwap", Chain: "BSC", Category: "DEX"},
	}

	result := FilterProtocols(input, "Ethereum", "")

	if len(result) != 1 {
		t.Fatalf("expected 1 protocol, got %d", len(result))
	}

	if result[0].Name != "Aave" {
		t.Fatalf("expected Aave, got %s", result[0].Name)
	}
}

func TestFilterProtocols_ByCategory(t *testing.T) {
	input := []models.ProtocolResponse{
		{Name: "Aave", Chain: "Ethereum", Category: "Lending"},
		{Name: "Uniswap", Chain: "Ethereum", Category: "DEX"},
	}

	result := FilterProtocols(input, "", "DEX")

	if len(result) != 1 {
		t.Fatalf("expected 1 protocol, got %d", len(result))
	}

	if result[0].Name != "Uniswap" {
		t.Fatalf("expected Uniswap, got %s", result[0].Name)
	}
}

func TestFilterProtocols_ByChainAndCategory(t *testing.T) {
	input := []models.ProtocolResponse{
		{Name: "Aave", Chain: "Ethereum", Category: "Lending"},
		{Name: "Aave", Chain: "Polygon", Category: "Lending"},
	}

	result := FilterProtocols(input, "Polygon", "Lending")

	if len(result) != 1 {
		t.Fatalf("expected 1 protocol, got %d", len(result))
	}

	if result[0].Chain != "Polygon" {
		t.Fatalf("expected Polygon chain, got %s", result[0].Chain)
	}
}