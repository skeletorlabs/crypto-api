package halving

import "time"

type State struct {
	CurrentBlock     int       `json:"currentBlock"`
	NextHalvingBlock int       `json:"nextHalvingBlock"`
	BlocksRemaining  int       `json:"blocksRemaining"`
	ProgressPercent  float64   `json:"progressPercent"`
	EstimatedDate    time.Time `json:"estimatedDate"`
	CurrentSubsidy   float64   `json:"currentSubsidy"`
	NextSubsidy      float64   `json:"nextSubsidy"`
}
