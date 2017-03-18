package analysis

import "time"

// Analyser - Performs analysis on proxy traffic
// and updates the statistics
type Analyser struct {
	MsgChannel chan string
	stats      *Statistics
}

// New - Returns a new Analyser Instance
func New() *Analyser {
	msgchan := make(chan string)
	return &Analyser{
		MsgChannel: msgchan,
		stats:      NewStatistics(),
	}
}

// Run - Starts the analyser
func (a *Analyser) Run() error {
	for {
		err := a.consume()
		if err != nil {
			return err
		}
	}
}

// Count - Provides the current count for messages of type stype
func (a *Analyser) Count(stype string) uint64 {
	return uint64(len(a.stats.count[stype]))
}

// TotalCount - Provides the current count for all messages
func (a *Analyser) TotalCount() uint64 {
	return uint64(
		len(a.stats.count["REQ"]) +
			len(a.stats.count["ACK"]) +
			len(a.stats.count["NAK"]))
}

func (a *Analyser) requestRate(timeUnit uint8) float32 {
	now := time.Now()
	num := 0
	for i := len(a.stats.count["REQ"]) - 1; i > 0; i-- {
		if now.After(a.stats.count["REQ"][i]) {
			num++
		}
	}
	return float32(num) / float32(timeUnit)
}

func (a *Analyser) responseRate(timeUnit uint8) float32 {
	return 0
}

func (a *Analyser) consume() error {
	msg := <-a.MsgChannel
	a.stats.Add(msg[:3])
	return nil
}
