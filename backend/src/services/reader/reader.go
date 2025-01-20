package reader

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
)

type Reader struct {
	repo repos.Repo
}

func (r *Reader) Services(ctx context.Context, query *ServicesArgs) ([]*Service, error) {
	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	services := []*Service{}

	if err := r.repo.ReadServerCalls(ctx, &services, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	errServices := []*Service{}

	if err := r.repo.ReadServerErrors(ctx, &errServices, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	serviceToErrCount := map[string]int{}

	for _, errService := range errServices {
		serviceToErrCount[errService.ServiceName] = errService.NumErrors
	}

	for _, service := range services {
		service.CallRate = float32(service.NumCalls) / float32(query.Period)

		if val, ok := serviceToErrCount[service.ServiceName]; ok {
			service.NumErrors = val
		}

		service.ErrorRate = float32(service.NumErrors) / float32(query.Period)
	}

	return services, nil
}

func (r *Reader) ServicesList(ctx context.Context) ([]string, error) {
	services := []string{}

	if err := r.repo.ReadServiceNames(ctx, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *Reader) ServiceDependencies(ctx context.Context, query *ServicesArgs) ([]*ServiceDependency, error) {
	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	serviceSpanDependencies := []*ServiceSpanDependency{}

	if err := r.repo.ReadSpanDependencies(ctx, &serviceSpanDependencies, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	spanToServiceName := map[string]string{}

	for _, dep := range serviceSpanDependencies {
		spanToServiceName[dep.SpanId] = dep.ServiceName
	}

	dependencies := map[string]*ServiceDependency{}

	for _, dep := range serviceSpanDependencies {
		parentToChild := spanToServiceName[dep.ParentSpanId] + "-" + spanToServiceName[dep.SpanId]

		if _, ok := dependencies[parentToChild]; !ok {
			dependencies[parentToChild] = &ServiceDependency{
				Parent:    spanToServiceName[dep.ParentSpanId],
				Child:     spanToServiceName[dep.SpanId],
				CallCount: 1,
			}
		} else {
			dependencies[parentToChild].CallCount += 1
		}
	}

	serviceDependencies := make([]*ServiceDependency, 0, len(dependencies))

	for _, dep := range dependencies {
		if len(dep.Parent) == 0 {
			continue
		}

		serviceDependencies = append(serviceDependencies, dep)
	}

	sort.Slice(serviceDependencies, func(i, j int) bool {
		return serviceDependencies[i].CallCount > serviceDependencies[j].CallCount
	})

	return serviceDependencies, nil
}

func (r *Reader) Operations(ctx context.Context, query *OperationsArgs) ([]string, error) {
	operations := []string{}

	if err := r.repo.ReadServiceSpecificOperations(ctx, &operations, query.ServiceName); err != nil {
		return nil, err
	}

	return operations, nil
}

func (r *Reader) Endpoints(ctx context.Context, query *EndpointsArgs) ([]*Endpoint, error) {
	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	topEndpoints := []*Endpoint{}

	if err := r.repo.ReadServiceSpecificEndpoints(ctx, &topEndpoints, query.ServiceName, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	return topEndpoints, nil
}

func (r *Reader) ServiceOverview(ctx context.Context, query *OverviewArgs) ([]*ServiceOverview, error) {
	interval := strconv.Itoa(int(query.StepSeconds / 60))

	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	serviceOverview := []*ServiceOverview{}

	if err := r.repo.ReadServiceSpecificServerCalls(ctx, &serviceOverview, query.ServiceName, interval, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	errServiceOverview := []*ServiceOverview{}

	if err := r.repo.ReadServiceSpecificServerErrors(ctx, &errServiceOverview, query.ServiceName, interval, startTimestamp, endTimestamp); err != nil {
		return nil, err
	}

	timeToErrCount := map[int64]int{}

	for _, errService := range errServiceOverview {
		t, _ := time.Parse(time.RFC3339Nano, errService.Time)

		timeToErrCount[int64(t.UnixNano())] = errService.NumErrors
	}

	result := []*ServiceOverview{}

	for _, service := range serviceOverview {
		overview := &ServiceOverview{
			Percentile50: service.Percentile50,
			Percentile95: service.Percentile95,
			Percentile99: service.Percentile99,
			NumCalls:     service.NumCalls,
		}

		overview.CallRate = float32(overview.NumCalls) / float32(query.StepSeconds)

		t, _ := time.Parse(time.RFC3339Nano, service.Time)

		overview.Timestamp = int64(t.UnixNano())

		if val, ok := timeToErrCount[overview.Timestamp]; ok {
			overview.NumErrors = val
		}

		overview.ErrorRate = float32(overview.NumErrors) / float32(query.StepSeconds)

		result = append(result, overview)
	}

	return result, nil
}

// https://opentelemetry.io/docs/specs/semconv/http/http-spans/
func (r *Reader) HttpOverview(ctx context.Context, query *OverviewArgs) ([]*HttpOverview, error) {
	return nil, nil
}

// https://opentelemetry.io/docs/specs/semconv/rpc/rpc-spans/
func (r *Reader) RpcOverview(ctx context.Context, query *OverviewArgs) ([]*RpcOverview, error) {
	return nil, nil
}

// https://opentelemetry.io/docs/specs/semconv/database/database-spans/
func (r *Reader) DBOverview(ctx context.Context, query *OverviewArgs) ([]*DBOverview, error) {
	return nil, nil
}

// https://opentelemetry.io/docs/specs/semconv/messaging/messaging-spans/
func (r *Reader) MessagingOverview(ctx context.Context, query *OverviewArgs) ([]*MessagingOverview, error) {
	return nil, nil
}

func (r *Reader) Tags(ctx context.Context, query *TagsArgs) ([]string, error) {
	tags := []string{}

	if err := r.repo.ReadServiceSpecificTags(ctx, &tags, query.ServiceName); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Reader) Spans(ctx context.Context, query *SpansArgs) ([]*SpanMatrix, error) {
	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	spans := []*Span{}

	if err := r.repo.ReadSpans(
		ctx,
		&spans,
		startTimestamp,
		endTimestamp,
		query.ServiceName,
		query.Name,
		query.Kind,
		query.MinDuration,
		query.MaxDuration,
		query.TagQueries...,
	); err != nil {
		return nil, err
	}

	spanMatrix := []*SpanMatrix{
		{
			Columns: []string{"Time", "SpanId", "ParentSpanId", "TraceId", "ServiceName", "Name", "Kind", "StatusCode", "Duration", "Tags"},
			Events:  make([][]interface{}, len(spans)),
		},
	}

	for i, span := range spans {
		event := span.ToEventValues()
		spanMatrix[0].Events[i] = event
	}

	return spanMatrix, nil
}

func (r *Reader) AggregatedSpans(ctx context.Context, query *AggregatedSpansArgs) ([]*AggregatedSpans, error) {
	interval := strconv.Itoa(int(query.StepSeconds / 60))

	startTimestamp := strconv.FormatInt(query.Start.UnixNano(), 10)

	endTimestamp := strconv.FormatInt(query.End.UnixNano(), 10)

	aggregationResults := []*AggregatedSpans{}

	if err := r.repo.ReadAggregatedSpans(
		ctx,
		&aggregationResults,
		query.Dimension,
		query.AggregationOption,
		interval,
		startTimestamp,
		endTimestamp,
		query.ServiceName,
		query.Name,
		query.Kind,
		query.MinDuration,
		query.MaxDuration,
		query.TagQueries...,
	); err != nil {
		return nil, err
	}

	result := []*AggregatedSpans{}

	for _, aggregatedSpan := range aggregationResults {
		agg := &AggregatedSpans{}

		t, _ := time.Parse(time.RFC3339Nano, aggregatedSpan.Time)

		agg.Timestamp = int64(t.UnixNano())
		agg.Value = aggregatedSpan.Value

		if query.AggregationOption == "rate_per_sec" {
			agg.Value = float32(aggregatedSpan.Value) / float32(query.StepSeconds)
		}

		result = append(result, agg)
	}

	return result, nil
}

func (r *Reader) SpansByTrace(ctx context.Context, query *SpansByTraceIdArgs) ([]*SpanMatrix, error) {
	spans := []*Span{}

	if err := r.repo.ReadTraceSpecificSpans(ctx, &spans, query.TraceId); err != nil {
		return nil, err
	}

	spanMatrix := []*SpanMatrix{
		{
			Columns: []string{"Time", "SpanId", "ParentSpanId", "TraceId", "ServiceName", "Name", "Kind", "StatusCode", "Duration", "Tags"},
			Events:  make([][]interface{}, len(spans)),
		},
	}

	for i, span := range spans {
		event := span.ToEventValues()
		spanMatrix[0].Events[i] = event
	}

	return spanMatrix, nil
}

func NewReader(repo repos.Repo) *Reader {
	r := &Reader{
		repo: repo,
	}

	return r
}
