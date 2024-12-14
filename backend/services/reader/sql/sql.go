package sql

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/services/reader"
)

type sqlReader struct {
	options reader.ReaderOptions
}

func (r *sqlReader) Options() reader.ReaderOptions {
	return r.options
}

func (r *sqlReader) Services(ctx context.Context, query *reader.ServicesArgs) ([]reader.Service, error) {
	return nil, nil
}

func (r *sqlReader) ServiceMapDependencies(ctx context.Context, query *reader.ServicesArgs) ([]reader.ServiceMapDependency, error) {
	return nil, nil
}

func (r *sqlReader) ServicesList(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (r *sqlReader) ServiceOverview(ctx context.Context, query *reader.ServiceOverviewArgs) ([]reader.ServiceOverview, error) {
	return nil, nil
}

func (r *sqlReader) ServiceDBOverview(ctx context.Context, query *reader.ServiceOverviewArgs) ([]reader.ServiceDBOverview, error) {
	return nil, nil
}

func (r *sqlReader) ServiceExternalAvgDuration(ctx context.Context, query *reader.ServiceOverviewArgs) ([]reader.ServiceExternalItem, error) {
	return nil, nil
}

func (r *sqlReader) ServiceExternalErrors(ctx context.Context, query *reader.ServiceOverviewArgs) ([]reader.ServiceExternalItem, error) {
	return nil, nil
}

func (r *sqlReader) ServiceExternal(ctx context.Context, query *reader.ServiceOverviewArgs) ([]reader.ServiceExternalItem, error) {
	return nil, nil
}

func (r *sqlReader) Operations(ctx context.Context, serviceName string) ([]string, error) {
	return nil, nil
}

func (r *sqlReader) TopEndpoints(ctx context.Context, query *reader.TopEndpointsArgs) ([]reader.TopEndpoints, error) {
	return nil, nil
}

func (r *sqlReader) Spans(ctx context.Context, query *reader.SpansArgs) ([]reader.Span, error) {
	return nil, nil
}

func (r *sqlReader) SpansAggregate(ctx context.Context, query *reader.SpansAggregateArgs) ([]reader.SpanAggregate, error) {
	return nil, nil
}

func (r *sqlReader) Tags(ctx context.Context, serviceName string) ([]reader.TagItem, error) {
	return nil, nil
}

func (r *sqlReader) Traces(ctx context.Context, traceId string) ([]reader.Span, error) {
	return nil, nil
}

func (r *sqlReader) Usage(ctx context.Context, query *reader.UsageArgs) ([]reader.Usage, error) {
	return nil, nil
}

func NewReader(opts ...reader.ReaderOption) reader.Reader {
	options := reader.NewReaderOptions(opts...)

	r := &sqlReader{
		options: options,
	}

	return r
}
