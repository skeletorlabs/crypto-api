package main

import (
	"log"
	"net/http"
	"os"

	"crypto-api/internal/cache"
	"crypto-api/internal/httpx"
	"crypto-api/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists, but don't panic if it's missing
	// as vars might be provided by the system/docker environment.
	_ = godotenv.Load()

	// Validate if essential variables are present (either from .env or system)
	if os.Getenv("FRED_API_KEY") == "" {
		log.Fatal("Critical Error: FRED_API_KEY is not set in environment")
	}

	log.Println("Environment variables loaded successfully")

	mux := http.NewServeMux()

	// Initialize caches
	marketCache := cache.NewMemoryCache()
	bitcoinCache := cache.NewMemoryCache()
	macroCache := cache.NewMemoryCache()
	intelligenceCache := cache.NewMemoryCache()

	v1 := http.NewServeMux()

	v1.HandleFunc("/health", httpx.HealthHandler)
	v1.HandleFunc("/price/", httpx.PriceHandler(marketCache))
	v1.HandleFunc("/chains", httpx.ChainsHandler(marketCache))
	v1.HandleFunc("/protocols", httpx.ProtocolsHandler(marketCache))
	v1.HandleFunc("/bitcoin/fees", httpx.BitcoinFeesHandler(bitcoinCache))
	v1.HandleFunc("/bitcoin/network", httpx.BitcoinNetworkHandler(bitcoinCache))
	v1.HandleFunc("/bitcoin/mempool", httpx.GetBitcoinMempoolHandler(bitcoinCache))
	v1.HandleFunc("/macro/liquidity", httpx.MacroHandler(macroCache))
	v1.HandleFunc("/bitcoin/valuation", httpx.ValuationHandler(marketCache, macroCache, intelligenceCache))
	v1.HandleFunc("/bitcoin/correlation", httpx.CorrelationHandler(marketCache, macroCache, intelligenceCache))

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	handler := middleware.Logging(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
