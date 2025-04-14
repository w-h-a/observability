package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/observability/backend/internal/clients/metrics"
	promrepo "github.com/w-h-a/observability/backend/internal/clients/metrics/prom"
	"github.com/w-h-a/observability/backend/internal/clients/traces"
	sqlrepo "github.com/w-h-a/observability/backend/internal/clients/traces/sql"
	"github.com/w-h-a/observability/backend/internal/config"
	httphandlers "github.com/w-h-a/observability/backend/internal/handlers/http"
	"github.com/w-h-a/observability/backend/internal/services/reader"
	"github.com/w-h-a/pkg/serverv2"
	httpserver "github.com/w-h-a/pkg/serverv2/http"
)

func Factory(tracesClient traces.Client, metricsClient metrics.Client) serverv2.Server {
	// clients
	sqlRepo := sqlrepo.NewRepo(
		traces.RepoWithClient(tracesClient),
		traces.RepoWithDatabase(config.TracesDB()),
		traces.RepoWithTable(config.TracesTable()),
	)

	promRepo := promrepo.NewRepo(
		metrics.RepoWithClient(metricsClient),
	)

	// services
	reader := reader.NewReader(sqlRepo, promRepo)

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
	httpMetrics := httphandlers.NewMetricsHandler(reader, httpRequestParser)

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
	router.Methods(http.MethodGet).Path("/api/v1/metrics").HandlerFunc(httpMetrics.GetMetrics)

	httpOpts := []serverv2.ServerOption{
		serverv2.ServerWithAddress(config.HttpAddress()),
		httpserver.HttpServerWithMiddleware(httphandlers.NewCORSMiddleware()),
	}

	httpOpts = append(httpOpts, opts...)

	httpServer := httpserver.NewServer(httpOpts...)

	httpServer.Handle(router)

	return httpServer
}
