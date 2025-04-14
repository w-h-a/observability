package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/w-h-a/observability/backend/src"
	"github.com/w-h-a/observability/backend/src/clients/traces"
	sqlrepo "github.com/w-h-a/observability/backend/src/clients/traces/sql"
	"github.com/w-h-a/observability/backend/src/config"
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
	repoClient := sqlrepo.NewClient(
		traces.ClientWithDriver(config.TracesStore()),
		traces.ClientWithAddrs(config.TracesStoreAddress()),
	)

	// server
	httpServer := src.Factory(repoClient)

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
