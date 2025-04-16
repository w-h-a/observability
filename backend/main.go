package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/w-h-a/observability/backend/internal"
	"github.com/w-h-a/observability/backend/internal/clients/metrics"
	promrepo "github.com/w-h-a/observability/backend/internal/clients/metrics/prom"
	"github.com/w-h-a/observability/backend/internal/clients/traces"
	sqlrepo "github.com/w-h-a/observability/backend/internal/clients/traces/sql"
	"github.com/w-h-a/observability/backend/internal/config"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
)

func main() {
	// config
	config.New()

	// resource
	name := fmt.Sprintf("%s.%s", config.Namespace(), config.Name())

	// log
	logBuffer := memoryutils.NewBuffer()

	logger := memorylog.NewLog(
		log.LogWithPrefix(name),
		memorylog.LogWithBuffer(logBuffer),
	)

	log.SetLogger(logger)

	// traces

	// clients
	tracesClient := sqlrepo.NewClient(
		traces.ClientWithDriver(config.TracesStore()),
		traces.ClientWithAddrs(config.TracesStoreAddress()),
	)

	metricsClient := promrepo.NewClient(
		metrics.ClientWithAddrs(config.MetricsStoreAddress()),
	)

	// server
	httpServer := internal.Factory(tracesClient, metricsClient)

	// wait group and error chan
	wg := &sync.WaitGroup{}
	errCh := make(chan error, 1)

	// start http server
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
