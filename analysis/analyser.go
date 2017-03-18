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
	since := time.Now().Add(-1 * time.Duration(timeUnit) * time.Second)
	count := a.eventsSince("REQ", since)

	return float32(count) / float32(timeUnit)
}

func (a *Analyser) eventsSince(eventType string, since time.Time) (num uint64) {
	num = 0
	for i := len(a.stats.count[eventType]) - 1; i >= 0; i-- {
		if a.stats.count[eventType][i].After(since) {
			num++
		}
	}
	return num
}

func (a *Analyser) responseRate(timeUnit uint8) float32 {
	since := time.Now().Add(-1 * time.Duration(timeUnit) * time.Second)
	countACK := a.eventsSince("ACK", since)
	countNAK := a.eventsSince("NAK", since)

	return float32(countACK+countNAK) / float32(timeUnit)
}

func (a *Analyser) consume() error {
	msg := <-a.MsgChannel
	a.stats.Add(msg[:3])
	return nil
}
