package reader

import "context"

type Reader interface {
	Services(ctx context.Context, query *ServicesArgs) ([]Service, error)
	ServiceMapDependencies(ctx context.Context, query *ServicesArgs) ([]ServiceMapDependency, error)
	ServicesList(ctx context.Context) ([]string, error)
	ServiceOverview(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceOverview, error)
	ServiceDBOverview(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceDBOverview, error)
	ServiceExternalAvgDuration(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error)
	ServiceExternalErrors(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error)
	ServiceExternal(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error)
	Operations(ctx context.Context, serviceName string) ([]string, error)
	TopEndpoints(ctx context.Context, query *TopEndpointsArgs) ([]TopEndpoints, error)
	Spans(ctx context.Context, query *SpansArgs) ([]Span, error)
	SpansAggregate(ctx context.Context, query *SpansAggregateArgs) ([]SpanAggregate, error)
	Tags(ctx context.Context, serviceName string) ([]TagItem, error)
	Traces(ctx context.Context, traceId string) ([]Span, error)
	Usage(ctx context.Context, query *UsageArgs) ([]Usage, error)
}
