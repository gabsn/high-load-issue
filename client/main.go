package main

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/dd-trace-go/tracer"
)

var (
	wg              = &sync.WaitGroup{}
	parallelConn    = 100
	traces          = getTestTrace(100, 100)
	defaultHostname = "localhost"
	defaultPort     = "8126"
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
		//createSpan()
		//sendTraces()
		encodeTraces()
		time.Sleep(500 * time.Microsecond)
	}
}

func createSpan() {
	span := tracer.NewRootSpan("test-go", "panic", "go")
	span.SetError(errors.New(strings.Repeat("a", 1000)))
	span.Finish()
}

func sendTraces() {
	transport := tracer.NewHTTPTransport(defaultHostname, defaultPort)
	_, err := transport.SendTraces(traces)
	if err != nil {
		println(err.Error())
	}
}

func encodeTraces() {
	encoder := tracer.NewMsgpackEncoder()
	err := encoder.EncodeTraces(traces)
	if err != nil {
		println(err.Error())
	}
}

// getTestTrace returns a list of traces that is composed by ``traceN`` number
// of traces, each one composed by ``size`` number of spans.
func getTestTrace(traceN, size int) [][]*tracer.Span {
	var traces [][]*tracer.Span
	for i := 0; i < traceN; i++ {
		trace := []*tracer.Span{}
		for j := 0; j < size; j++ {
			trace = append(trace, getTestSpan())
		}
		traces = append(traces, trace)
	}
	return traces
}

func getTestSpan() *tracer.Span {
	return &tracer.Span{
		TraceID:  42,
		SpanID:   52,
		ParentID: 42,
		Type:     "web",
		Service:  "high.throughput",
		Name:     "sending.events",
		Resource: "SEND /data",
		Start:    1481215590883401105,
		Duration: 1000000000,
		Meta:     map[string]string{"http.host": "192.168.0.1"},
		Metrics:  map[string]float64{"http.monitor": 41.99},
	}
}
