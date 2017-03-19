// +build integration

package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/antoniou/go-foo-proxy/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProxyTestSuite struct {
	suite.Suite
}

type mockServer struct {
	mock.Mock
}

func (m *mockServer) Run() {
	l, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Mock Server Listening on localhost:8001")
	buff := make([]byte, 0xffff)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		count, _ := conn.Read(buff)
		b := buff[:count]
		// Handle connections in a new goroutine.
		m.handleRequest(string(b))
		conn.Write([]byte("ACK 1 Hey\n"))
	}
}

func (m *mockServer) handleRequest(data string) {
	m.Called(data)
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Run() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8002")
	fmt.Print("Messsage to send: ")
	text := "REQ 1 Hello"
	// send to socket
	fmt.Fprintf(conn, text+"\n")

	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	m.handleReply(message)
}

func (m *mockClient) handleReply(message string) {
	fmt.Print("Message from server: " + message)
	m.Called(message)
}

func (suite *ProxyTestSuite) TestProxyRun() {
	m := &mockServer{}
	c := &mockClient{}
	m.On("handleRequest", "REQ 1 Hello\n").Return(nil)
	c.On("handleReply", "ACK 1 Hey\n").Return(nil)

	go m.Run()
	time.Sleep(1 * time.Second)
	localAddr := ":8002"
	remoteAddr := "localhost:8001"
	p := proxy.New(localAddr, remoteAddr)
	go p.Run()
	c.Run()
	time.Sleep(1 * time.Second)

	// Expect Server to receive message
	m.AssertExpectations(suite.T())
	m.AssertCalled(suite.T(), "handleRequest", "REQ 1 Hello\n")

	// Expect reponse from the client
	c.AssertExpectations(suite.T())
	c.AssertCalled(suite.T(), "handleReply", "ACK 1 Hey\n")

	// Send SIGUSR1 signal to Proxy
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(1 * time.Second)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	assert.Contains(suite.T(), buf.String(), "msg_total: 2")
	assert.Contains(suite.T(), buf.String(), "msg_req: 1")
	assert.Contains(suite.T(), buf.String(), "msg_ack: 1")
	assert.Contains(suite.T(), buf.String(), "msg_nak: 0")
	assert.Contains(suite.T(), buf.String(), "request_rate_1s: 0.000000")
	assert.Contains(suite.T(), buf.String(), "request_rate_10s: 0.1")
	assert.Contains(suite.T(), buf.String(), "response_rate_1s: 0.000000")
	assert.Contains(suite.T(), buf.String(), "response_rate_10s: 0.1")

}

func TestProxyTestSuite(t *testing.T) {
	suite.Run(t, new(ProxyTestSuite))
}
