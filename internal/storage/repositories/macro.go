package repositories

import (
	"context"
	"crypto-api/internal/engine/bitcoin/correlation"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MacroRepository struct {
	pool *pgxpool.Pool
}

func NewMacroRepository(pool *pgxpool.Pool) *MacroRepository {
	return &MacroRepository{pool: pool}
}

// SaveM2 persists a macro M2 data point using the source-provided date.
func (r *MacroRepository) SaveM2(ctx context.Context, supply float64, sourceDate time.Time) error {
	query := `
        INSERT INTO macro_stats (m2_supply, source_date, created_at) 
        VALUES ($1, $2, $3)
    `
	_, err := r.pool.Exec(ctx, query, supply, sourceDate, time.Now().UTC())
	return err
}

// GetLatestM2 returns the most recent macro value by source date.
func (r *MacroRepository) GetLatestM2(ctx context.Context) (float64, time.Time, error) {
	var supply float64
	var sourceDate time.Time

	query := `
        SELECT m2_supply, source_date 
        FROM macro_stats 
        ORDER BY source_date DESC 
        LIMIT 1
    `

	err := r.pool.QueryRow(ctx, query).Scan(&supply, &sourceDate)
	return supply, sourceDate, err
}

// GetM2History returns M2 historical data in chronological order.
func (r *MacroRepository) GetM2History(ctx context.Context, limit int) ([]correlation.DataPoint, error) {
	query := `
        SELECT m2_supply, source_date 
        FROM macro_stats 
        ORDER BY source_date DESC 
        LIMIT $1
    `

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []correlation.DataPoint

	for rows.Next() {
		var p correlation.DataPoint
		if err := rows.Scan(&p.Value, &p.Date); err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	// Reverse to chronological order (oldest → newest)
	for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
		points[i], points[j] = points[j], points[i]
	}

	return points, nil
}
