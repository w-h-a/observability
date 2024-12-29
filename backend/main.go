package main

import (
	"sync"
	"time"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
	sqlrepo "github.com/w-h-a/trace-blame/backend/src/clients/repos/sql"
	"github.com/w-h-a/trace-blame/backend/src/config"
)

func main() {
	// config
	config.NewConfig()

	// clients
	repoClient := sqlrepo.NewClient(
		repos.ClientWithDriver(config.Store()),
		repos.ClientWithAddrs(config.StoreAddress()),
	)

	// server
	httpServer := src.ServerFactory(repoClient)

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
