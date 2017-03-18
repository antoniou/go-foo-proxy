package analysis

import "time"

// Statistics - A struct holding partially processed
// analysis data
type Statistics struct {
	count map[string][]time.Time
}

// NewStatistics - Initialises a new Statistics struct
func NewStatistics() *Statistics {
	count := make(map[string][]time.Time)
	count["REQ"] = make([]time.Time, 0, 1000000)
	count["ACK"] = make([]time.Time, 0, 1000000)
	count["NAK"] = make([]time.Time, 0, 1000000)
	return &Statistics{
		count: count,
	}
}

func (s *Statistics) Add(stype string) error {
	s.count[stype] = append(s.count[stype], time.Now())
	return nil
}
