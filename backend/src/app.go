package src

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"github.com/w-h-a/pkg/telemetry/log"
	memorylog "github.com/w-h-a/pkg/telemetry/log/memory"
	"github.com/w-h-a/pkg/utils/memoryutils"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
	memorystore "github.com/w-h-a/trace-blame/backend/src/clients/store/memory"
	sqlstore "github.com/w-h-a/trace-blame/backend/src/clients/store/sql"
	"github.com/w-h-a/trace-blame/backend/src/config"
	httphandlers "github.com/w-h-a/trace-blame/backend/src/handlers/http"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/repos"
	sqlrepo "github.com/w-h-a/trace-blame/backend/src/repos/sql"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

var (
	defaultStoreClients = map[string]func(...store.ClientOption) store.Client{
		"clickhouse": sqlstore.NewClient,
		"memory":     memorystore.NewClient,
	}
)

func AppFactory(storeClient store.Client) serverv2.Server {
	// config
	config.NewConfig()

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
	sqlRepo := sqlrepo.NewRepo(
		repos.RepoWithClient(getStoreClient(storeClient)),
		repos.RepoWithDatabase(config.DB()),
		repos.RepoWithTable(config.Table()),
	)

	// services
	reader := reader.NewReader(sqlRepo)

	// base server options
	opts := []serverv2.ServerOption{
		serverv2.ServerWithNamespace(config.Namespace()),
		serverv2.ServerWithName(config.Name()),
		serverv2.ServerWithVersion(config.Version()),
	}

	// create http server
	router := mux.NewRouter()

	requestParser := &utils.RequestParser{}

	httpServices := httphandlers.NewServicesHandler(reader, requestParser)

	router.Methods(http.MethodGet).Path("/api/v1/services").HandlerFunc(httpServices.GetServices)
	router.Methods(http.MethodGet).Path("/api/v1/services/list").HandlerFunc(httpServices.GetServicesList)
	router.Methods(http.MethodGet).Path("/api/v1/services/dependencies").HandlerFunc(httpServices.GetServiceDependencies)
	router.Methods(http.MethodGet).Path("/api/v1/services/operations").HandlerFunc(httpServices.GetOperations)
	router.Methods(http.MethodGet).Path("/api/v1/services/overview").HandlerFunc(httpServices.GetServiceOverview)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}

func getStoreClient(storeClient store.Client) store.Client {
	if storeClient != nil {
		return storeClient
	}

	storeClientBuilder, exists := defaultStoreClients[config.Store()]
	if !exists && len(config.Store()) > 0 {
		log.Fatalf("store %s not supported", config.Store())
	} else if !exists {
		return memorystore.NewClient()
	}

	return storeClientBuilder(
		store.ClientWithDriver(config.Store()),
		store.ClientWithAddrs(config.StoreAddress()),
	)
}
