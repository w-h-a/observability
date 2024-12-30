package src

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
	sqlrepo "github.com/w-h-a/trace-blame/backend/src/clients/repos/sql"
	"github.com/w-h-a/trace-blame/backend/src/config"
	httphandlers "github.com/w-h-a/trace-blame/backend/src/handlers/http"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func ServerFactory(repoClient repos.Client) serverv2.Server {
	// clients
	sqlRepo := sqlrepo.NewRepo(
		repos.RepoWithClient(repoClient),
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
	httpService := httphandlers.NewServiceHandler(reader, requestParser)
	httpSpans := httphandlers.NewSpansHandler(reader, requestParser)

	router.Methods(http.MethodGet).Path("/api/v1/services").HandlerFunc(httpServices.GetServices)
	router.Methods(http.MethodGet).Path("/api/v1/services/list").HandlerFunc(httpServices.GetServicesList)
	router.Methods(http.MethodGet).Path("/api/v1/services/dependencies").HandlerFunc(httpServices.GetServiceDependencies)
	router.Methods(http.MethodGet).Path("/api/v1/service/operations").HandlerFunc(httpService.GetOperations)
	router.Methods(http.MethodGet).Path("/api/v1/service/endpoints").HandlerFunc(httpService.GetEndpoints)
	router.Methods(http.MethodGet).Path("/api/v1/service/overview").HandlerFunc(httpService.GetServiceOverview)
	router.Methods(http.MethodGet).Path("/api/v1/spans/trace").HandlerFunc(httpSpans.GetSpansByTraceId)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}
