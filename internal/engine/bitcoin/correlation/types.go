package correlation

import (
	"time"

	"crypto-api/internal/models"
)

// Result contains the computed Pearson correlation
// along with dataset metadata.
type Result struct {
	models.Meta

	Coefficient float64   `json:"coefficient"`
	SampleCount int       `json:"sample_count"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// DataPoint represents a time-series value
// used for alignment before correlation computation.
type DataPoint struct {
	Date  time.Time
	Value float64
}
