package storage

import (
	"context"
	"crypto-api/internal/models"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	Pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context) (*PostgresStore, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{Pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.Pool.Close()
}

// SaveSnapshot persists a structural intelligence snapshot.
func (s *PostgresStore) SaveSnapshot(ctx context.Context, snap models.IntelligenceSnapshot) error {
	query := `
    INSERT INTO bitcoin_intelligence_snapshots 
    (
      price_usd, m2_supply_billions, btc_m2_ratio, pearson_correlation, 
      block_height, hashrate_ths, difficulty, avg_block_time, 
      network_health_score, trend_status, source_attribution
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.Pool.Exec(ctx, query,
		snap.PriceUSD,
		snap.M2SupplyBillions,
		snap.BTCM2Ratio,
		snap.Correlation,
		snap.BlockHeight,
		snap.HashrateTHs,
		snap.Difficulty,
		snap.AvgBlockTime,
		snap.NetworkHealthScore,
		snap.TrendStatus,
		snap.SourceAttribution,
	)

	return err
}

// GetLatestSnapshot returns the most recent intelligence snapshot.
func (s *PostgresStore) GetLatestSnapshot(ctx context.Context) (*models.IntelligenceSnapshot, error) {
	query := `
    SELECT 
      id, created_at, price_usd, m2_supply_billions, btc_m2_ratio, 
      pearson_correlation, block_height, hashrate_ths, difficulty, 
      avg_block_time, network_health_score, trend_status, source_attribution
    FROM bitcoin_intelligence_snapshots 
    ORDER BY created_at DESC 
    LIMIT 1`

	var snap models.IntelligenceSnapshot

	err := s.Pool.QueryRow(ctx, query).Scan(
		&snap.ID,
		&snap.CreatedAt,
		&snap.PriceUSD,
		&snap.M2SupplyBillions,
		&snap.BTCM2Ratio,
		&snap.Correlation,
		&snap.BlockHeight,
		&snap.HashrateTHs,
		&snap.Difficulty,
		&snap.AvgBlockTime,
		&snap.NetworkHealthScore,
		&snap.TrendStatus,
		&snap.SourceAttribution,
	)

	if err != nil {
		return nil, err
	}

	return &snap, nil
}
