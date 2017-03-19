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

## Overview



## Future Work/Improvements:
1. Do not open a new TCP connection to server on every incoming connection
2. More efficient analysis of data
3. Logging
