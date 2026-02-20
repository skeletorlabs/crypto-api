package halving

import "time"

// State represents the current Bitcoin halving cycle status
// derived from block height and average block time.
type State struct {
	CurrentBlock     int       `json:"currentBlock"`
	NextHalvingBlock int       `json:"nextHalvingBlock"`
	BlocksRemaining  int       `json:"blocksRemaining"`
	ProgressPercent  float64   `json:"progressPercent"`
	EstimatedDate    time.Time `json:"estimatedDate"`
	CurrentSubsidy   float64   `json:"currentSubsidy"`
	NextSubsidy      float64   `json:"nextSubsidy"`
}
