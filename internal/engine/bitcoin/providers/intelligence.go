package providers

import (
	"context"
	"fmt"
	"log"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/config"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
)

const assetSymbol = "BTC"

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
	now := time.Now().UTC()

	latestSnapshot, _ := p.intelligenceRepo.GetLatest(ctx)

	// --- PRICE ---
	var price float64
	source := "cache"

	data, found := cache.Get[models.IntelligencePrice](
		p.intelligenceCache,
		cache.KeyIntelligencePrice(assetSymbol),
	)

	if !found {
		log.Printf("[intelligence] Price missing in cache, attempting database fallback...")

		if latestSnapshot == nil {
			return fmt.Errorf("BTC price missing in cache and no database fallback available")
		}

		price = latestSnapshot.PriceUSD
		source = "database_fallback"
	} else {
		price = data.Price
	}

	// Persist daily price before snapshot generation
	today := now.Truncate(24 * time.Hour)
	if err := p.priceHistoryRepo.SavePrice(ctx, models.PriceHistory{
		Timestamp: today,
		Asset:     assetSymbol,
		PriceUSD:  price,
		Source:    source,
	}); err != nil {
		log.Printf("[intelligence] Failed to save price history: %v", err)
	}

	// --- NETWORK ---
	var netData models.BitcoinNetworkResponse
	cachedNet, found := cache.Get[models.BitcoinNetworkResponse](
		p.intelligenceCache,
		cache.KeyBitcoinNetwork,
	)

	if found {
		netData = cachedNet
	} else {
		dbNet, err := p.networkRepo.GetLatest(ctx)
		if err != nil {
			return fmt.Errorf("network data missing for snapshot: %w", err)
		}
		netData = *dbNet
	}

	// --- MACRO ---
	m2Supply, _, err := p.macroRepo.GetLatestM2(ctx)
	if err != nil {
		return fmt.Errorf("macro data missing for snapshot: %w", err)
	}

	// --- CORRELATION ---
	pearsonValue := 0.0

	m2History, err := p.macroRepo.GetM2History(ctx, config.CorrelationLookbackDays)
	if err != nil {
		log.Printf("[intelligence] M2 history unavailable: %v", err)
	}
	priceSeries, err := p.priceHistoryRepo.GetPriceSeries(ctx, assetSymbol, config.CorrelationLookbackDays)
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

	if res, errCorr := correlation.Compute(m2History, btcHistory); errCorr == nil {
		pearsonValue = res.Coefficient
	} else {
		log.Printf("[intelligence] Correlation skipped: %v", errCorr)
	}

	// --- VALUATION ---
	vState := valuation.Compute(price, m2Supply)
	healthScore := valuation.CalculateNetworkHealth(netData.AvgBlockTimeSeconds)
	hasPrev := latestSnapshot != nil
	prevAvg := 0.0
	if hasPrev {
		prevAvg = latestSnapshot.AvgBlockTime
	}

	trendStatus := valuation.CalculateTrend(
		netData.AvgBlockTimeSeconds,
		prevAvg,
		hasPrev,
	)

	// --- SNAPSHOT ---
	snapshot := models.IntelligenceSnapshot{
		SnapshotDate:       today,
		CreatedAt:          now,
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
