package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/antoniou/go-foo-proxy/analysis"
)

var connid = uint64(0)

// Proxy - Listens to incoming tcp connections and
// forwards data between local address laddr and
// remote address remoteAddr
type Proxy struct {
	laddr, raddr *net.TCPAddr
	lconn, rconn io.ReadWriteCloser
	analyser     *analysis.Analyser
	reporter     *analysis.Reporter
}

// FooMessage - A Foo protocol message
type FooMessage struct {
	Type string
	Seq  int
	Data string // TODO: Make buff
	raw  []byte
}

// New - Creates a Proxy instance
func New(localAddr string, remoteAddr string) *Proxy {
	laddr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		fmt.Printf("Could not resolve local address: %s", err)
		os.Exit(1)
	}

	raddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("Failed to resolve remote address: %s", err)
		os.Exit(1)
	}

	analyser := analysis.New()
	return &Proxy{
		laddr:    laddr,
		raddr:    raddr,
		analyser: analyser,
		reporter: analysis.NewReporter(analyser),
	}
}

// Run - The proxy starts listening for connections and
// handles them
func (p *Proxy) Run() {
	listener, err := net.ListenTCP("tcp", p.laddr)
	if err != nil {
		fmt.Printf("Failed to open local port to listen: %s", err)
		os.Exit(1)
	}

	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		fmt.Printf("Remote connection failed: %s", err)
		return
	}
	defer p.rconn.Close()

	go p.analyser.Run()
	go p.reporter.Run()

	for {
		p.lconn, err = listener.AcceptTCP()

		if err != nil {
			fmt.Printf("Failed to accept connection '%s'", err)
			continue
		}
		connid++
		go p.pipe(p.lconn, p.rconn)
		go p.pipe(p.rconn, p.lconn)
	}
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	for {
		msg, err := p.readMessage(src)
		if err != nil {
			return
		}
		p.analyser.MsgChannel <- string(msg.raw)
		p.writeMessage(msg, dst)
	}
}

// readMessage - Reads from source connection and
// tries to construct a Message
func (p *Proxy) readMessage(src io.ReadWriter) (*FooMessage, error) {
	messageComplete := false
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)

	pos := 0
	for messageComplete != true {
		n, err := src.Read(buff[pos:])
		if err != nil {
			return nil, err
		}
		pos += n
		fmt.Printf("Read total data %s", buff[:pos])
		if bytes.ContainsAny(buff[:pos], "\n") {
			messageComplete = true
		}
	}

	return &FooMessage{
		raw: buff[:pos],
	}, nil
}

func (p *Proxy) writeMessage(msg *FooMessage, dst io.ReadWriter) error {
	fmt.Printf("Sending data %s", msg.raw)
	_, err := dst.Write(msg.raw)
	if err != nil {
		fmt.Printf("Write failed '%s'\n", err)
		return err
	}
	return nil
}
