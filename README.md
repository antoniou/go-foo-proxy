# Go-Foo-Proxy [![Build Status](https://travis-ci.org/antoniou/go-foo-proxy.svg?branch=master)](https://travis-ci.org/antoniou/go-foo-proxy)

This is an intercepting proxy for a Foo Protocol, written in Go.

## Installation
To install go-foo-proxy, you'll need to have Golang installed and environment variable [$GOPATH appropriately set](https://golang.org/doc/install).
```bash
$ go get github.com/antoniou/go-foo-proxy
```

## Usage
There needs to be a server for the Foo protocol already running. Assuming that the server is running on localhost:8001, you can start the proxy like so:

```bash
$ go-foo-proxy -listen=:8002 -forward localhost:8001
Proxying from :8002 to localhost:8001
```

To output the statistics, send a SIGUSR1 signal to foo-proxy:
```bash
$ kill -SIGUSR1 $(pidof go-foo-proxy)
```

## Overview
The main components of the solution are:
  * Proxy: The intercepting proxy that connects to a remote server and accepts connections from clients. It consequently pipes data from/to the client and server. On every incoming request, it calls the Analyser to update the collected analysis data.
  * Analyser: The analyser is responsible for updating the Statistics (See next item) and for calculating the requested metrics.
  * Statistics: Collection of data that will be used for analysis. At the moment, the main data structure is count, which holds 3 lists, one for each of "REQ","ACK", "NAK". Each list contains sorted timestamps of events of the specific type. The analyser will place new data in Statistics every time a new Request/Response is being made.
  * Reporter: The reporter awaits for a SIGUSR1 signal. When the Reporter is notified by such a signal, it queries the Analyser for the metrics and provides it with the requested format to Stdout


## Future Work/Improvements:
1. *Proxy/Server Connection Pooling*: At the moment, everytime the proxy accepts a new client connection it also establishes a new server connection. That results in a lot of connections to the server being maintained. To resolve this, the Proxy needs to be modified so that
connections to the server are reused.
2. *Better Logging*: At the moment we are not using a logger. It would be more appropriate to have a separate logger that outputs events with a verbosity level
3. *Better exception handling*: Because of the limited time spent on this assignment, special situations (e.g, incorrect message format) are not being handled correctly.
4. *Performance optimisations*: Data Analysis has not been optimised. As a result the solution is not expected to scale to large number of requests.  
