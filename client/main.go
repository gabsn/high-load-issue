package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

var (
	wg                   = &sync.WaitGroup{}
	currentRequest int64 = 0
	parallelConn         = 1000
	addr                 = "http://localhost:8080"
)

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
		doRequest(curr)
	}
}

func doRequest(i int64) {
	resp, err := http.Get(addr)
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
