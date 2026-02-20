package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypto-api/internal/app"
	"crypto-api/internal/engine/intelligence"
)

func main() {

	// Graceful shutdown context (SIGINT, SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
	defer application.Store.Close()

	log.Println("[worker] Starting background workers...")

	intelligence.StartMarketWorker(ctx, application.Caches.Intelligence)
	intelligence.StartMacroWorker(ctx, application.Caches.Macro, application.Repos.Macro)
	intelligence.StartNetworkWorker(ctx, application.Caches.Bitcoin, application.Repos.Network)
	intelligence.StartDailyWorker(ctx, application.Provider)

	<-ctx.Done()
	log.Println("[worker] Shutting down...")
}
