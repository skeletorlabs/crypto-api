package intelligence

import (
	"context"
	"crypto-api/internal/models"
)

// Snapshottable represents any asset engine capable of producing
// an IntelligenceSnapshot.
type Snapshottable interface {
	GetSnapshot(ctx context.Context) (models.IntelligenceSnapshot, error)
}
