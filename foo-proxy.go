package main

import (
	"flag"
	"fmt"

	"github.com/antoniou/go-foo-proxy/proxy"
)

var (
	localAddr  = flag.String("listen", ":8002", "local address")
	remoteAddr = flag.String("forward", "localhost:8001", "remote address")
	verbose    = flag.Bool("v", false, "display server actions")
)

func main() {
	flag.Parse()

	fmt.Printf("Proxying from %v to %v\n", *localAddr, *remoteAddr)

	p := proxy.New(*localAddr, *remoteAddr)
	p.Run()

}
