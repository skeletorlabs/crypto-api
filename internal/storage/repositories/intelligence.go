package repositories

import (
	"context"
	"crypto-api/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IntelligenceRepository struct {
	pool *pgxpool.Pool
}

func NewIntelligenceRepository(pool *pgxpool.Pool) *IntelligenceRepository {
	return &IntelligenceRepository{pool: pool}
}

// SaveSnapshot persists a new intelligence snapshot.
func (r *IntelligenceRepository) SaveSnapshot(ctx context.Context, s models.IntelligenceSnapshot) error {
	if r.pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}

	query := `
		INSERT INTO intelligence_snapshots (
			snapshot_date,
			created_at,
			price_usd,
			m2_supply,
			btc_m2_ratio, 
			correlation,
			block_height,
			hashrate_ths,
			difficulty,
			avg_block_time,
			network_health_score,
			trend_status,
			source_attribution
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		)
		ON CONFLICT (snapshot_date) DO NOTHING
		`

	_, err := r.pool.Exec(ctx, query,
		s.SnapshotDate,       // $1
		s.CreatedAt,          // $2
		s.PriceUSD,           // $3
		s.M2SupplyBillions,   // $4
		s.BTCM2Ratio,         // $5
		s.Correlation,        // $6
		s.BlockHeight,        // $7
		s.HashrateTHs,        // $8
		s.Difficulty,         // $9
		s.AvgBlockTime,       // $10
		s.NetworkHealthScore, // $11
		s.TrendStatus,        // $12
		s.SourceAttribution,  // $13
	)

	return err
}

// GetLatest returns the most recent snapshot.
func (r *IntelligenceRepository) GetLatest(ctx context.Context) (*models.IntelligenceSnapshot, error) {
	if r.pool == nil {
		return nil, fmt.Errorf("database pool is not initialized")
	}

	var s models.IntelligenceSnapshot

	query := `
    SELECT 
      id,
      snapshot_date,
      created_at,
      price_usd,
      m2_supply,
      btc_m2_ratio,
      correlation, 
      block_height,
      hashrate_ths,
      difficulty,
      avg_block_time, 
      network_health_score,
      trend_status,
      source_attribution 
    FROM intelligence_snapshots 
    ORDER BY snapshot_date DESC 
    LIMIT 1`

	err := r.pool.QueryRow(ctx, query).Scan(
		&s.ID,
		&s.SnapshotDate,
		&s.CreatedAt,
		&s.PriceUSD,
		&s.M2SupplyBillions,
		&s.BTCM2Ratio,
		&s.Correlation,
		&s.BlockHeight,
		&s.HashrateTHs,
		&s.Difficulty,
		&s.AvgBlockTime,
		&s.NetworkHealthScore,
		&s.TrendStatus,
		&s.SourceAttribution,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
