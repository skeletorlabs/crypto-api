package repositories

import (
	"context"
	"crypto-api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PriceHistoryRepository struct {
	Pool *pgxpool.Pool
}

func NewPriceHistoryRepository(pool *pgxpool.Pool) *PriceHistoryRepository {
	return &PriceHistoryRepository{Pool: pool}
}

// SavePrice persists a price record, updating on conflict.
func (r *PriceHistoryRepository) SavePrice(ctx context.Context, p models.PriceHistory) error {
	query := `
		INSERT INTO price_history (timestamp, asset, price_usd, source)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (timestamp, asset) DO UPDATE 
		SET price_usd = EXCLUDED.price_usd,
		    source = EXCLUDED.source;
	`

	_, err := r.Pool.Exec(ctx, query, p.Timestamp, p.Asset, p.PriceUSD, p.Source)
	return err
}

// GetPriceSeries returns the latest N price points in chronological order.
func (r *PriceHistoryRepository) GetPriceSeries(
	ctx context.Context,
	asset string,
	limit int,
) ([]models.PriceHistory, error) {

	query := `
		SELECT timestamp, price_usd
		FROM price_history
		WHERE asset = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.Pool.Query(ctx, query, asset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var series []models.PriceHistory

	for rows.Next() {
		var ph models.PriceHistory
		if err := rows.Scan(&ph.Timestamp, &ph.PriceUSD); err != nil {
			return nil, err
		}
		series = append(series, ph)
	}

	// Reverse to chronological order (oldest → newest)
	for i, j := 0, len(series)-1; i < j; i, j = i+1, j-1 {
		series[i], series[j] = series[j], series[i]
	}

	return series, nil
}
