package analysis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStats struct {
	mock.Mock
}

func (m *mockStats) Count(stype string) []time.Time {
	return []time.Time{
		time.Now().Add(-6 * time.Second),
		time.Now().Add(-3 * time.Second),
	}
}

func (m *mockStats) Add(stype string) error {
	m.Called(stype)
	return nil
}

func NewTestAnalyser(m *mockStats) *Analyser {
	msgchan := make(chan string)
	return &Analyser{
		MsgChannel: msgchan,
		stats:      m,
	}
}

type AnalyserTestSuite struct {
	suite.Suite
}

func (suite *AnalyserTestSuite) TestCount() {
	m := &mockStats{}
	a := NewTestAnalyser(m)

	assert.Equal(suite.T(), uint64(2), a.Count("REQ"))
	assert.Equal(suite.T(), uint64(2), a.Count("ACK"))
	assert.Equal(suite.T(), uint64(2), a.Count("NAK"))
	assert.Equal(suite.T(), uint64(6), a.TotalCount())
}

func (suite *AnalyserTestSuite) TestRequestRate() {
	m := &mockStats{}
	a := NewTestAnalyser(m)

	assert.Equal(suite.T(), float32(0), a.RequestRate(1))
	assert.Equal(suite.T(), float32(0.2), a.RequestRate(5))
	assert.Equal(suite.T(), float32(0.2), a.RequestRate(10))

	assert.Equal(suite.T(), float32(0), a.ResponseRate(1))
	assert.Equal(suite.T(), float32(0.4), a.ResponseRate(5))
	assert.Equal(suite.T(), float32(0.4), a.ResponseRate(10))
}

func (suite *AnalyserTestSuite) TestRunAnalyser() {
	m := &mockStats{}
	a := NewTestAnalyser(m)
	m.On("Add", "REQ").Return(nil)
	go func() {
		a.MsgChannel <- "REQ 1 Hello\n"
	}()
	a.consume()
	m.Add("REQ")
	m.AssertExpectations(suite.T())
	m.AssertCalled(suite.T(), "Add", "REQ")
}

func TestAnalyserTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyserTestSuite))
}
