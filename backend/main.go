package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"github.com/w-h-a/trace-blame/backend/config"
)

func main() {
	name := fmt.Sprintf("%s.%s", config.Namespace, config.Name)

	// logger
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix(name),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	// tracer

	// clients

	// services

	// base server opts
	opts := []serverv2.ServerOption{
		serverv2.ServerWithNamespace(config.Namespace),
		serverv2.ServerWithName(config.Name),
		serverv2.ServerWithVersion(config.Version),
	}

	// create http server
	router := mux.NewRouter()

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	// wait group and error chan
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)

	// run http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh <- httpServer.Start()
	}()

	// block
	err := <-errCh
	if err != nil {
		log.Errorf("failed to start server: %+v", err)
	}

	// graceful shutdown
	wait := make(chan struct{})

	go func() {
		defer close(wait)
		wg.Wait()
	}()

	select {
	case <-wait:
	case <-time.After(30 * time.Second):
	}

	log.Info("successfully stopped server")
}
