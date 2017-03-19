package analysis

import (
	"fmt"
	"time"
)

// Reporter - Listens to SIGUSR1 signal and
// provides a report on the proxy traffic
type Reporter struct {
	analyser AnalyserI
}

// NewReporter - Creates a new Reporter struct
func NewReporter(a *Analyser) *Reporter {
	return &Reporter{
		analyser: a,
	}
}

// Report - Creates and returns a statistics report
func (r *Reporter) Report() string {
	return fmt.Sprintf("msg_total: %d\n", r.analyser.TotalCount()) +
		fmt.Sprintf("msg_req: %d\n", r.analyser.Count("REQ")) +
		fmt.Sprintf("msg_ack: %d\n", r.analyser.Count("ACK")) +
		fmt.Sprintf("msg_nak: %d\n", r.analyser.Count("NAK")) +
		fmt.Sprintf("request_rate_1s: %f\n", r.analyser.RequestRate(1)) +
		fmt.Sprintf("request_rate_10s: %f\n", r.analyser.RequestRate(10)) +
		fmt.Sprintf("response_rate_1s: %f\n", r.analyser.ResponseRate(1)) +
		fmt.Sprintf("response_rate_10s: %f\n", r.analyser.ResponseRate(10))
}

// Run - Starts the reporter
func (r *Reporter) Run() error {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println(r.Report())
	}
}
