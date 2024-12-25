package reader

import (
	"context"
	"sort"
	"strconv"

	"github.com/w-h-a/trace-blame/backend/src/repos"
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

func (r *Reader) ServicesList(ctx context.Context) ([]string, error) {
	// TODO: log/trace?

	services := []string{}

	if err := r.repo.ReadDistinctServiceNames(ctx, &services); err != nil {
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
