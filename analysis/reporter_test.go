package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockAnalyser struct {
	mock.Mock
}

func (m *mockAnalyser) Count(stype string) uint64 {
	return 100
}

func (m *mockAnalyser) TotalCount() uint64 {
	return 1000
}

func (m *mockAnalyser) RequestRate(timeUnit uint8) float32 {
	return 10.1
}

func (m *mockAnalyser) ResponseRate(timeUnit uint8) float32 {
	return 10.2
}

func NewTestReporter() *Reporter {
	return &Reporter{
		analyser: &mockAnalyser{},
	}
}

type ReporterTestSuite struct {
	suite.Suite
}

func (suite *ReporterTestSuite) TestReport() {
	r := NewTestReporter()
	report := r.Report()
	assert.Contains(suite.T(), report, "msg_total: 1000")
	assert.Contains(suite.T(), report, "msg_req: 100")
	assert.Contains(suite.T(), report, "msg_ack: 100")
	assert.Contains(suite.T(), report, "msg_nak: 100")
	assert.Contains(suite.T(), report, "request_rate_1s: 10.1")
	assert.Contains(suite.T(), report, "request_rate_10s: 10.1")
	assert.Contains(suite.T(), report, "response_rate_1s: 10.2")
	assert.Contains(suite.T(), report, "response_rate_10s: 10.2")
}

func TestReporterTestSuite(t *testing.T) {
	suite.Run(t, new(ReporterTestSuite))
}
