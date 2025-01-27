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
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func AppFactory(repoClient repos.Client) serverv2.Server {
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

	httpRequestParser := &httphandlers.RequestParser{}

	httpServices := httphandlers.NewServicesHandler(reader, httpRequestParser)
	httpService := httphandlers.NewServiceHandler(reader, httpRequestParser)
	httpTraces := httphandlers.NewTracesHandler(reader, httpRequestParser)
	httpSpans := httphandlers.NewSpansHandler(reader, httpRequestParser)

	router.Methods(http.MethodGet).Path("/api/v1/services").HandlerFunc(httpServices.GetServices)
	router.Methods(http.MethodGet).Path("/api/v1/services/list").HandlerFunc(httpServices.GetServicesList)
	router.Methods(http.MethodGet).Path("/api/v1/services/dependencies").HandlerFunc(httpServices.GetServiceDependencies)
	router.Methods(http.MethodGet).Path("/api/v1/service/tags").HandlerFunc(httpService.GetTags)
	router.Methods(http.MethodGet).Path("/api/v1/service/operations").HandlerFunc(httpService.GetOperations)
	router.Methods(http.MethodGet).Path("/api/v1/service/endpoints").HandlerFunc(httpService.GetEndpoints)
	router.Methods(http.MethodGet).Path("/api/v1/service/overview").HandlerFunc(httpService.GetServiceOverview)
	router.Methods(http.MethodGet).Path("/api/v1/traces").HandlerFunc(httpTraces.GetTraces)
	router.Methods(http.MethodGet).Path("/api/v1/spans").HandlerFunc(httpSpans.GetSpans)
	router.Methods(http.MethodGet).Path("/api/v1/spans/aggregated").HandlerFunc(httpSpans.GetAggregatedSpans)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
		httpserver.HttpServerWithMiddleware(httphandlers.NewCORSMiddleware()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}
