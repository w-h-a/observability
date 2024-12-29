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
	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
	memoryrepo "github.com/w-h-a/trace-blame/backend/src/clients/repos/memory"
	sqlrepo "github.com/w-h-a/trace-blame/backend/src/clients/repos/sql"
	"github.com/w-h-a/trace-blame/backend/src/config"
	httphandlers "github.com/w-h-a/trace-blame/backend/src/handlers/http"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

var (
	defaultRepoClients = map[string]func(...repos.ClientOption) repos.Client{
		"clickhouse": sqlrepo.NewClient,
		"memory":     memoryrepo.NewClient,
	}
)

func AppFactory(repoClient repos.Client) serverv2.Server {
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
		repos.RepoWithClient(getRepoClient(repoClient)),
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
	router.Methods(http.MethodGet).Path("/api/v1/services/endpoints").HandlerFunc(httpServices.GetEndpoints)
	router.Methods(http.MethodGet).Path("/api/v1/services/overview").HandlerFunc(httpServices.GetServiceOverview)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}

func getRepoClient(repoClient repos.Client) repos.Client {
	if repoClient != nil {
		return repoClient
	}

	repoClientBuilder, exists := defaultRepoClients[config.Store()]
	if !exists && len(config.Store()) > 0 {
		log.Fatalf("store %s not supported", config.Store())
	} else if !exists {
		return memoryrepo.NewClient()
	}

	return repoClientBuilder(
		repos.ClientWithDriver(config.Store()),
		repos.ClientWithAddrs(config.StoreAddress()),
	)
}
