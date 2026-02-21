package main

import (
	"context"
	"log"
	stdhttp "net/http"
	"time"

	"crypto-api/internal/app"

	"crypto-api/internal/http/bitcoin"
	"crypto-api/internal/http/chains"
	"crypto-api/internal/http/health"
	"crypto-api/internal/http/intelligence"
	"crypto-api/internal/http/macro"
	"crypto-api/internal/http/market"
	"crypto-api/internal/http/protocols"
)

func main() {

	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
	defer application.Store.Close()

	mux := stdhttp.NewServeMux()
	v1 := stdhttp.NewServeMux()

	// ===========================================================================
	// Health
	// ===========================================================================
	// GET /health
	// Liveness probe. Confirms API and DB connectivity.
	v1.HandleFunc("/health",
		health.HealthHandler(application.Store),
	)

	// ===========================================================================
	// Market (Real-Time Layer)
	// External real-time data from exchanges and aggregators.
	// ===========================================================================

	// GET /market/price/{token}
	// Real-time spot price from external providers (Binance/Coinbase/Kraken fallback).
	v1.HandleFunc("/market/price/",
		market.MarketPriceHandler(application.Caches.Market),
	)

	// GET /market/protocols
	// DeFi protocol list with TVL and metadata (DefiLlama).
	v1.HandleFunc("/market/protocols",
		protocols.ProtocolsHandler(application.Caches.Market),
	)

	// GET /market/chains
	// Supported blockchain list and aggregated metrics.
	v1.HandleFunc("/market/chains",
		chains.ChainsHandler(application.Caches.Market),
	)

	// ===========================================================================
	// Intelligence (Structural Layer)
	// Engine-derived state synchronized with macro + network metrics.
	// ===========================================================================

	// GET /intelligence/price/{token}
	// Structural price used by the intelligence engine (snapshot-aligned).
	v1.HandleFunc("/intelligence/price/",
		intelligence.IntelligencePriceHandler(
			application.Caches.Intelligence,
			application.Repos.Intelligence,
		),
	)

	// GET /intelligence/snapshot
	// Full structural state (price, M2 ratio, correlation, trend, health).
	v1.HandleFunc("/intelligence/snapshot",
		intelligence.IntelligenceSnapshotHandler(
			application.Caches.Intelligence,
			application.Repos.Intelligence,
		),
	)

	// ===========================================================================
	// Bitcoin (On-Chain Metrics)
	// Raw + derived Bitcoin network data.
	// ===========================================================================

	// GET /bitcoin/network
	// Current Bitcoin network state with halving + trend computation.
	v1.HandleFunc("/bitcoin/network",
		bitcoin.BitcoinNetworkHandler(
			application.Caches.Bitcoin,
			application.Repos.Network,
		),
	)

	// GET /bitcoin/correlation
	// BTC vs liquidity historical correlation series.
	v1.HandleFunc("/bitcoin/correlation",
		bitcoin.BitcoinCorrelationHandler(
			application.Repos.Intelligence,
		),
	)

	// GET /bitcoin/fees
	// Current mempool fee estimation tiers.
	v1.HandleFunc("/bitcoin/fees",
		bitcoin.BitcoinFeesHandler(application.Caches.Bitcoin),
	)

	// GET /bitcoin/mempool
	// Current mempool size and fee pressure metrics.
	v1.HandleFunc("/bitcoin/mempool",
		bitcoin.BitcoinMempoolHandler(application.Caches.Bitcoin),
	)

	// ===========================================================================
	// Macro (Liquidity Layer)
	// Global liquidity inputs used in structural calculations.
	// ===========================================================================

	// GET /macro/m2
	// Latest M2 money supply data used by intelligence engine.
	v1.HandleFunc("/macro/m2",
		macro.MacroHandler(
			application.Caches.Macro,
			application.Repos.Macro,
		),
	)

	mux.Handle("/v1/", stdhttp.StripPrefix("/v1", v1))

	server := &stdhttp.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("http server listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
