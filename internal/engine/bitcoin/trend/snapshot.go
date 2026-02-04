package trend

import "time"

type Snapshot struct {
	Timestamp       time.Time
	MempoolVSize    int64
	AvgBlockTimeSec float64
}
