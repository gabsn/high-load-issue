package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	wg                   = &sync.WaitGroup{}
	currentRequest int64 = 0
	totalRequests  int64 = 100000
	parallelConn         = 1000
	addr                 = "http://localhost:8080"
	dialer               = &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport = &http.Transport{
		MaxIdleConns:        4096,
		MaxIdleConnsPerHost: 2048,
		DialContext: func(ctx context.Context, network, addr string) (c net.Conn, err error) {
			c, err = dialer.DialContext(ctx, network, addr)
			return
		},
	}
	client = &http.Client{
		Transport: transport,
	}
)

func init() {
	flag.StringVar(&addr, "addr", addr, "default http://localhost:8080 target addr")
	flag.IntVar(&parallelConn, "p", parallelConn, "default 100 parallel connections")
	flag.Int64Var(&totalRequests, "t", totalRequests, "default 10000 requests")
	flag.Parse()
}

func main() {
	wg.Add(parallelConn)
	for i := 0; i < parallelConn; i++ {
		go requestRoutine()
	}
	wg.Wait()
}

func requestRoutine() {
	defer wg.Done()
	for {
		curr := atomic.AddInt64(&currentRequest, 1)
		if curr > totalRequests {
			break
		}
		doRequest(curr)
	}
	transport.CloseIdleConnections()
}

func doRequest(i int64) {
	req, _ := http.NewRequest("GET", addr, nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error on request: %d (%s)\n", i, err)
		os.Exit(1)
	}
	if _, err = io.Copy(ioutil.Discard, resp.Body); err != nil {
		fmt.Printf("error on body read: %d (%s)\n", i, err)
		os.Exit(1)
	}
	resp.Body.Close()
}
