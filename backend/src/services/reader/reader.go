package reader

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/src/repos"
)

type Reader struct {
	repo repos.Repo
}

func (r *Reader) Services(ctx context.Context, query *ServicesArgs) ([]Service, error) {
	return nil, nil
}

func (r *Reader) ServiceMapDependencies(ctx context.Context, query *ServicesArgs) ([]ServiceMapDependency, error) {
	return nil, nil
}

func (r *Reader) ServicesList(ctx context.Context) ([]string, error) {
	// TODO: log/trace?

	services := []string{}

	if err := r.repo.ReadServicesList(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *Reader) ServiceOverview(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceOverview, error) {
	return nil, nil
}

func (r *Reader) ServiceDBOverview(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceDBOverview, error) {
	return nil, nil
}

func (r *Reader) ServiceExternalAvgDuration(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error) {
	return nil, nil
}

func (r *Reader) ServiceExternalErrors(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error) {
	return nil, nil
}

func (r *Reader) ServiceExternal(ctx context.Context, query *ServiceOverviewArgs) ([]ServiceExternalItem, error) {
	return nil, nil
}

func (r *Reader) Operations(ctx context.Context, serviceName string) ([]string, error) {
	return nil, nil
}

func (r *Reader) TopEndpoints(ctx context.Context, query *TopEndpointsArgs) ([]TopEndpoints, error) {
	return nil, nil
}

func (r *Reader) Spans(ctx context.Context, query *SpansArgs) ([]Span, error) {
	return nil, nil
}

func (r *Reader) SpansAggregate(ctx context.Context, query *SpansAggregateArgs) ([]SpanAggregate, error) {
	return nil, nil
}

func (r *Reader) Tags(ctx context.Context, serviceName string) ([]TagItem, error) {
	return nil, nil
}

func (r *Reader) Traces(ctx context.Context, traceId string) ([]Span, error) {
	return nil, nil
}

func (r *Reader) Usage(ctx context.Context, query *UsageArgs) ([]Usage, error) {
	return nil, nil
}

func NewReader(repo repos.Repo) *Reader {
	r := &Reader{
		repo: repo,
	}

	return r
}
