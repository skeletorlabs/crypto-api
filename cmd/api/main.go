package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"crypto-api/internal/app"
	"crypto-api/internal/httpx"
)

func main() {

	ctx := context.Background()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("❌ Fatal: %v", err)
	}
	defer application.Store.Close()

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	// Infra
	v1.HandleFunc("/health", httpx.HealthHandler(application.Store))

	// Market
	v1.HandleFunc("/market/price/",
		httpx.MarketPriceHandler(application.Caches.Market))

	// Intelligence
	v1.HandleFunc("/intelligence/price/",
		httpx.IntelligencePriceHandler(
			application.Caches.Intelligence,
			application.Repos.Intelligence,
		))

	v1.HandleFunc("/bitcoin/intelligence",
		httpx.IntelligenceHandler(
			application.Caches.Intelligence,
			application.Repos.Intelligence,
		))

	v1.HandleFunc("/bitcoin/network",
		httpx.BitcoinNetworkHandler(
			application.Caches.Bitcoin,
			application.Repos.Network,
		))

	v1.HandleFunc("/macro/liquidity",
		httpx.MacroHandler(
			application.Caches.Macro,
			application.Repos.Macro,
		))

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🎸 Skeletor API running on :8080")
	log.Fatal(server.ListenAndServe())
}
