package main

import (
	"sync"
	"strings"
	"errors"

	"github.com/DataDog/dd-trace-go/tracer"
)

var (
	wg           = &sync.WaitGroup{}
	parallelConn = 1000
	addr         = "http://localhost:8080"
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
		doRequest()
	}
}

func doRequest() {
	span := tracer.NewRootSpan("test-go", "panic", "go")
	span.SetError(errors.New(strings.Repeat("a", 10)))
	span.Finish()
}
