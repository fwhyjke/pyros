package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"git.ecom.tech/ecom/dev/pap/backend/go/log"
	"github.com/fwhyjke/pyros/pyroscope"
)

func main() {
	log := log.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	profiler := pyroscope.New(
		pyroscope.WithPyroscopeServerAddress("http://localhost:17420"),
		pyroscope.WithServiceName("test-app"),
	)
	profiler.WithLogger(log)
	profiler.Start(context.TODO())

	go leakMem()

	<-ctx.Done()
	profiler.Stop(context.TODO())
}

var leak [][]byte

func leakMem() {
	for {
		chunk := make([]byte, 1024*1000)
		leak = append(leak, chunk)
		time.Sleep(time.Second)
	}
}
