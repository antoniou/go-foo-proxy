package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Reporter - Listens to SIGUSR1 signal and
// provides a report on the proxy traffic
type Reporter struct {
	analyser AnalyserI
}

// Report - A statistics report
type Report struct {
	MsgTotal        uint64  `json:"msg_total"`
	MsgReq          uint64  `json:"msg_req"`
	MsgAck          uint64  `json:"msg_ack"`
	MsgNak          uint64  `json:"msg_nak"`
	RequestRate1s   float32 `json:"request_rate_1s"`
	RequestRate10s  float32 `json:"request_rate_10s"`
	ResponseRate1s  float32 `json:"response_rate_1s"`
	ResponseRate10s float32 `json:"response_rate_10s"`
}

// NewReporter - Creates a new Reporter struct
func NewReporter(a *Analyser) *Reporter {
	return &Reporter{
		analyser: a,
	}
}

// Report - Creates and returns a statistics report
func (r *Reporter) Report() string {
	mreport := &Report{
		MsgTotal:        r.analyser.TotalCount(),
		MsgReq:          r.analyser.Count("REQ"),
		MsgAck:          r.analyser.Count("ACK"),
		MsgNak:          r.analyser.Count("NAK"),
		RequestRate1s:   r.analyser.RequestRate(1),
		RequestRate10s:  r.analyser.RequestRate(10),
		ResponseRate1s:  r.analyser.ResponseRate(1),
		ResponseRate10s: r.analyser.ResponseRate(10),
	}

	b, err := json.Marshal(*mreport)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

// Run - Starts the reporter
func (r *Reporter) Run() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for {
		s := <-c
		switch s {
		case syscall.SIGUSR1:
			fmt.Println(r.Report())
		default:
			return fmt.Errorf("Got signal %s", s)
		}
	}
}
