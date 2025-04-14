package prom

import (
	"context"
	"time"

	"github.com/w-h-a/observability/backend/internal/clients/metrics"
	"github.com/w-h-a/pkg/telemetry/log"
)

type promRepo struct {
	options metrics.RepoOptions
}

func (r *promRepo) Options() metrics.RepoOptions {
	return r.options
}

func (r *promRepo) ReadMetrics(ctx context.Context, dest interface{}, dimension string, step time.Duration, start, end time.Time, opts ...metrics.QueryOption) error {
	queryOptions := metrics.NewQueryOptions(opts...)

	queryFactory, ok := Queries[dimension]
	if !ok {
		return metrics.ErrUnsupportedQuery
	}

	query := queryFactory(queryOptions.ServiceName)

	if err := r.options.Client.Read(ctx, dest, query, start, end, step); err != nil {
		log.Errorf("metric client failed to read: %v", err)
		return metrics.ErrProcessingQuery
	}

	return nil
}

func NewRepo(opts ...metrics.RepoOption) metrics.Repo {
	options := metrics.NewRepoOptions(opts...)

	r := &promRepo{
		options: options,
	}

	return r
}
