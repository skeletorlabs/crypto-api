package main

import (
	"log"
	"net/http"
	"os"

	"crypto-api/internal/cache"
	"crypto-api/internal/httpx"
	"crypto-api/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	priceCache := cache.NewMemoryCache()
	chainsCache := cache.NewMemoryCache()
	protocolsCache := cache.NewMemoryCache()
	feesCache := cache.NewMemoryCache()
	networkCache := cache.NewMemoryCache()
	mempoolCache := cache.NewMemoryCache()

	v1 := http.NewServeMux()

	v1.HandleFunc("/health", httpx.HealthHandler)
	v1.HandleFunc("/price/", httpx.PriceHandler(priceCache))
	v1.HandleFunc("/chains", httpx.ChainsHandler(chainsCache))
	v1.HandleFunc("/protocols", httpx.ProtocolsHandler((protocolsCache)))
	v1.HandleFunc("/bitcoin/fees", httpx.BitcoinFeesHandler(feesCache))
	v1.HandleFunc("/bitcoin/network", httpx.BitcoinNetworkHandler(networkCache))
	v1.HandleFunc("/bitcoin/mempool", httpx.GetBitcoinMempoolHandler(mempoolCache))

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	handler := middleware.Logging(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
