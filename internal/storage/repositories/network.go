package repositories

import (
	"context"
	"crypto-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NetworkRepository struct {
	pool *pgxpool.Pool
}

func NewNetworkRepository(pool *pgxpool.Pool) *NetworkRepository {
	return &NetworkRepository{pool: pool}
}

// Save persists network statistics.
func (r *NetworkRepository) Save(ctx context.Context, n models.BitcoinNetworkResponse) error {
	query := `
        INSERT INTO network_stats 
        (block_height, hashrate_ths, avg_block_time, difficulty, created_at) 
        VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx,
		query,
		n.BlockHeight,
		n.HashrateTHs,
		n.AvgBlockTimeSeconds,
		n.Difficulty,
		time.Now().UTC(),
	)

	return err
}

// GetLatest returns the most recent network snapshot.
func (r *NetworkRepository) GetLatest(ctx context.Context) (*models.BitcoinNetworkResponse, error) {
	var n models.BitcoinNetworkResponse

	query := `
        SELECT block_height, hashrate_ths, avg_block_time, difficulty, created_at 
        FROM network_stats 
        ORDER BY created_at DESC 
        LIMIT 1`

	err := r.pool.QueryRow(ctx, query).Scan(
		&n.BlockHeight,
		&n.HashrateTHs,
		&n.AvgBlockTimeSeconds,
		&n.Difficulty,
		&n.Meta.UpdatedAt,
	)

	n.Meta.Cached = false

	return &n, err
}
