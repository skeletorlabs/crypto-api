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

func (r *NetworkRepository) DeleteOlderThan(ctx context.Context, days int) error {
	query := `
        DELETE FROM network_stats
        WHERE created_at < NOW() - ($1 || ' days')::interval
    `
	_, err := r.pool.Exec(ctx, query, days)
	return err
}

// GetPrevious returns the second most recent network snapshot.
func (r *NetworkRepository) GetPrevious(ctx context.Context) (*models.BitcoinNetworkResponse, error) {
	var n models.BitcoinNetworkResponse

	query := `
        SELECT block_height, hashrate_ths, avg_block_time, difficulty, created_at
        FROM network_stats
        ORDER BY created_at DESC
        OFFSET 1
        LIMIT 1
    `

	err := r.pool.QueryRow(ctx, query).Scan(
		&n.BlockHeight,
		&n.HashrateTHs,
		&n.AvgBlockTimeSeconds,
		&n.Difficulty,
		&n.Meta.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	n.Meta.Cached = false
	return &n, nil
}
