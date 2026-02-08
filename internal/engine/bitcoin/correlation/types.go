package correlation

import (
	"crypto-api/internal/models"
	"time"
)

// CorrelationResult holds the final calculation and metadata
type Result struct {
	models.Meta
	Coefficient float64   `json:"coefficient"`
	SampleCount int       `json:"sample_count"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// DataPoint is a helper for aligning dates before computation
type DataPoint struct {
	Date  time.Time
	Value float64
}
