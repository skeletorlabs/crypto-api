package providers

import (
	"context"
	"fmt"
	"log"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
)

type IntelligenceProvider struct {
	intelligenceCache *cache.MemoryCache
	macroCache        *cache.MemoryCache
	networkRepo       *repositories.NetworkRepository
	macroRepo         *repositories.MacroRepository
	intelligenceRepo  *repositories.IntelligenceRepository
	priceHistoryRepo  *repositories.PriceHistoryRepository
}

func NewIntelligenceProvider(
	intelCache, macroCache *cache.MemoryCache,
	networkRepo *repositories.NetworkRepository,
	macroRepo *repositories.MacroRepository,
	intelligenceRepo *repositories.IntelligenceRepository,
	priceHistoryRepo *repositories.PriceHistoryRepository,
) *IntelligenceProvider {
	return &IntelligenceProvider{
		intelligenceCache: intelCache,
		macroCache:        macroCache,
		networkRepo:       networkRepo,
		macroRepo:         macroRepo,
		intelligenceRepo:  intelligenceRepo,
		priceHistoryRepo:  priceHistoryRepo,
	}
}

// GenerateFullSnapshot builds and persists the daily intelligence snapshot.
func (p *IntelligenceProvider) GenerateFullSnapshot(ctx context.Context) error {
	// Prevent duplicate snapshot for the same UTC day
	if latest, err := p.intelligenceRepo.GetLatest(ctx); err == nil && latest != nil {
		if latest.CreatedAt.UTC().Format("2006-01-02") ==
			time.Now().UTC().Format("2006-01-02") {
			log.Printf("[intelligence] Snapshot for today already exists. Skipping.")
			return nil
		}
	}

	// --- PRICE ---
	var price float64
	source := "cache"

	data, found := cache.Get[models.IntelligencePrice](
		p.intelligenceCache,
		cache.KeyIntelligencePrice("BTC"),
	)

	if !found {
		log.Printf("[intelligence] Price missing in cache, attempting database fallback...")
		latest, err := p.intelligenceRepo.GetLatest(ctx)
		if err != nil || latest == nil {
			return fmt.Errorf("BTC price missing in cache and no database fallback available")
		}
		price = latest.PriceUSD
		source = "database_fallback"
	} else {
		price = data.Price
	}

	// Persist daily price before snapshot generation
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if err := p.priceHistoryRepo.SavePrice(ctx, models.PriceHistory{
		Timestamp: today,
		Asset:     "BTC",
		PriceUSD:  price,
		Source:    source,
	}); err != nil {
		log.Printf("[intelligence] Failed to save price history: %v", err)
	}

	// --- NETWORK ---
	netData, err := p.networkRepo.GetLatest(ctx)
	if err != nil {
		return fmt.Errorf("network data missing for snapshot: %w", err)
	}

	// --- MACRO ---
	m2Supply, _, err := p.macroRepo.GetLatestM2(ctx)
	if err != nil {
		return fmt.Errorf("macro data missing for snapshot: %w", err)
	}

	// --- CORRELATION ---
	pearsonValue := 0.0

	m2History, _ := p.macroRepo.GetM2History(ctx, 30)
	priceSeries, err := p.priceHistoryRepo.GetPriceSeries(ctx, "BTC", 30)
	if err != nil {
		log.Printf("[intelligence] Error fetching internal history: %v", err)
	}

	var btcHistory []correlation.DataPoint
	for _, ph := range priceSeries {
		btcHistory = append(btcHistory, correlation.DataPoint{
			Date:  ph.Timestamp,
			Value: ph.PriceUSD,
		})
	}

	// Forward-fill alignment
	var alignedM2 []correlation.DataPoint
	var finalBTCHistory []correlation.DataPoint
	m2Index := 0
	lastValue := 0.0
	hasStarted := false

	for _, btc := range btcHistory {
		for m2Index < len(m2History) &&
			(m2History[m2Index].Date.Before(btc.Date) ||
				m2History[m2Index].Date.Equal(btc.Date)) {

			lastValue = m2History[m2Index].Value
			m2Index++
			hasStarted = true
		}

		if hasStarted {
			alignedM2 = append(alignedM2, correlation.DataPoint{
				Date:  btc.Date,
				Value: lastValue,
			})
			finalBTCHistory = append(finalBTCHistory, btc)
		}
	}

	if len(alignedM2) >= 7 && len(alignedM2) == len(finalBTCHistory) {
		if res, errCorr := correlation.Compute(alignedM2, finalBTCHistory); errCorr == nil {
			pearsonValue = res.Coefficient
		}
	} else {
		log.Printf(
			"[intelligence] Insufficient or mismatched data (M2: %d, BTC: %d). Skipping Pearson.",
			len(alignedM2),
			len(finalBTCHistory),
		)
	}

	// --- VALUATION ---
	vState := valuation.Compute(price, m2Supply)
	healthScore := valuation.CalculateNetworkHealth(netData.AvgBlockTimeSeconds)

	prev, _ := p.intelligenceRepo.GetLatest(ctx)
	trendStatus := valuation.CalculateTrend(netData.AvgBlockTimeSeconds, prev)

	// --- SNAPSHOT ---
	snapshot := models.IntelligenceSnapshot{
		CreatedAt:          time.Now().UTC(),
		PriceUSD:           price,
		M2SupplyBillions:   m2Supply,
		BTCM2Ratio:         vState.Ratio,
		Correlation:        pearsonValue,
		BlockHeight:        netData.BlockHeight,
		HashrateTHs:        netData.HashrateTHs,
		Difficulty:         netData.Difficulty,
		NetworkHealthScore: healthScore,
		TrendStatus:        trendStatus,
		AvgBlockTime:       netData.AvgBlockTimeSeconds,
		SourceAttribution:  vState.Description,
	}

	if err := p.intelligenceRepo.SaveSnapshot(ctx, snapshot); err != nil {
		return err
	}

	cache.Set(
		p.intelligenceCache,
		cache.KeyIntelligenceLatestSnapshot,
		snapshot,
		cache.TTLIntelligenceSnapshot,
	)

	return nil
}

// Hydrate warms caches from persisted data.
func (p *IntelligenceProvider) Hydrate(ctx context.Context) error {
	log.Println("[intelligence] Hydrating caches from repositories...")

	if net, err := p.networkRepo.GetLatest(ctx); err == nil && net != nil {
		cache.Set(p.intelligenceCache, cache.KeyBitcoinNetwork, *net, cache.TTLNetworkStats)
	}

	if supply, lastUpdate, err := p.macroRepo.GetLatestM2(ctx); err == nil {
		resp := models.MacroResponse{
			Meta: models.Meta{
				UpdatedAt: lastUpdate,
				Cached:    false,
			},
			M2Supply: models.M2Details{
				Value:    supply,
				Unit:     "Billions",
				DateTime: lastUpdate,
			},
		}
		cache.Set(p.macroCache, cache.KeyMacroM2Supply, resp, cache.TTLMacroData)
	}

	if snap, err := p.intelligenceRepo.GetLatest(ctx); err == nil && snap != nil {
		cache.Set(p.intelligenceCache, cache.KeyIntelligenceLatestSnapshot, *snap, cache.TTLIntelligenceSnapshot)
	}

	return nil
}
