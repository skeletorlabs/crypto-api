package app

import (
	"context"
	"log"

	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/providers"
	"crypto-api/internal/storage"
	"crypto-api/internal/storage/repositories"

	"github.com/joho/godotenv"
)

type Caches struct {
	Market       *cache.MemoryCache
	Bitcoin      *cache.MemoryCache
	Macro        *cache.MemoryCache
	Intelligence *cache.MemoryCache
}

type Repositories struct {
	Network      *repositories.NetworkRepository
	Macro        *repositories.MacroRepository
	Intelligence *repositories.IntelligenceRepository
	PriceHistory *repositories.PriceHistoryRepository
}

type App struct {
	Store    *storage.PostgresStore
	Caches   Caches
	Repos    Repositories
	Provider *providers.IntelligenceProvider
}

func New(ctx context.Context) (*App, error) {

	// Load environment variables from .env if present
	_ = godotenv.Load()

	// Initialize database connection
	store, err := storage.NewPostgresStore(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize in-memory caches
	caches := Caches{
		Market:       cache.NewMemoryCache(),
		Bitcoin:      cache.NewMemoryCache(),
		Macro:        cache.NewMemoryCache(),
		Intelligence: cache.NewMemoryCache(),
	}

	// Initialize repositories
	repos := Repositories{
		Network:      repositories.NewNetworkRepository(store.Pool),
		Macro:        repositories.NewMacroRepository(store.Pool),
		Intelligence: repositories.NewIntelligenceRepository(store.Pool),
		PriceHistory: repositories.NewPriceHistoryRepository(store.Pool),
	}

	// Initialize intelligence provider (orchestrator)
	provider := providers.NewIntelligenceProvider(
		caches.Intelligence,
		caches.Macro,
		repos.Network,
		repos.Macro,
		repos.Intelligence,
		repos.PriceHistory,
	)

	// Hydrate caches at startup
	log.Println("[bootstrap] Hydrating caches...")
	if err := provider.Hydrate(ctx); err != nil {
		log.Printf("Hydration warning: %v", err)
	}

	return &App{
		Store:    store,
		Caches:   caches,
		Repos:    repos,
		Provider: provider,
	}, nil
}
