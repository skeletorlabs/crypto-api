package cache

import "fmt"

// --- Market Keys (Dados Voláteis/Brutos) ---
const (
	KeyMarketChains    = "market:chains:list"
	KeyMarketProtocols = "market:protocols:list"
)

func KeyMarketPrice(token string) string {
	return fmt.Sprintf("market:price:%s", token)
}

func KeyMarketHistory(token string) string {
	return fmt.Sprintf("market:history:%s", token)
}

// --- Bitcoin/Network Keys (Estado da Rede) ---
const (
	KeyBitcoinFees    = "bitcoin:fees:status"
	KeyBitcoinNetwork = "bitcoin:network:status" // Usado para hidratar o NetworkResponse
	KeyBitcoinMempool = "bitcoin:mempool:status"
)

// --- Macro Keys (Dados Econômicos) ---
const (
	KeyMacroM2Supply  = "macro:m2:supply"
	KeyMacroM2History = "macro:m2:history"
)

// --- Intelligence Keys (O Cérebro) ---
// Note que aqui separamos o SNAPSHOT (o todo) dos componentes (partes)
const (
	// KeyIntelligenceLatestSnapshot é a chave mestre para o endpoint /intelligence
	// Ela deve guardar o models.IntelligenceSnapshot completo.
	KeyIntelligenceLatestSnapshot = "intelligence:snapshot:latest"

	// Estas podem ser usadas para caches específicos de cálculos intermediários
	KeyIntelligenceValuation   = "intelligence:valuation:bitcoin"
	KeyIntelligenceCorrelation = "intelligence:correlation:bitcoin_m2"
)

// Para o caso de querermos guardar preços históricos processados pela inteligência
func KeyIntelligencePrice(token string) string {
	return fmt.Sprintf("intelligence:price:%s", token)
}
