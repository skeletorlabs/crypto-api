package halving

import (
	"testing"
)

func TestCompute_StartOfEpoch(t *testing.T) {
	// Bloco 840.000 é exatamente o início de um novo Halving (Abril/2024)
	currentBlock := 840000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	// O próximo halving deve ser 840.000 + 210.000
	if res.NextHalvingBlock != 1050000 {
		t.Errorf("expected next halving at 1050000, got %d", res.NextHalvingBlock)
	}

	// O progresso deve ser 0% no primeiro bloco
	if res.ProgressPercent != 0 {
		t.Errorf("expected progress 0, got %f", res.ProgressPercent)
	}
}

func TestCompute_MidEpoch(t *testing.T) {
	// Bloco 945.000 é exatamente o meio do caminho entre 840k e 1050k
	currentBlock := 945000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	// O progresso deve ser exatamente 50%
	if res.ProgressPercent != 50.0 {
		t.Errorf("expected progress 50.0, got %f", res.ProgressPercent)
	}

	// Devem faltar 105.000 blocos
	if res.BlocksRemaining != 105000 {
		t.Errorf("expected 105000 blocks remaining, got %d", res.BlocksRemaining)
	}
}

func TestCompute_LastBlockOfEpoch(t *testing.T) {
	// Bloco 1.049.999 é o último bloco antes do próximo halving
	currentBlock := 1049999
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	if res.BlocksRemaining != 1 {
		t.Errorf("expected 1 block remaining, got %d", res.BlocksRemaining)
	}

	// O progresso deve estar quase em 100%
	if res.ProgressPercent < 99.9 {
		t.Errorf("expected progress near 100, got %f", res.ProgressPercent)
	}
}
