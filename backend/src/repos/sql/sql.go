package sql

import (
	"context"
	"fmt"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/src/repos"
)

type sqlRepo struct {
	options repos.RepoOptions
}

func (r *sqlRepo) Options() repos.RepoOptions {
	return r.options
}

func (r *sqlRepo) ReadServerCalls(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`, r.options.Database, r.options.Table, startTimestamp, endTimestamp)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("store client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServerErrors(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT ServiceName as serviceName, count(*) as numErrors FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' AND StatusCode='Error' GROUP BY serviceName`, r.options.Database, r.options.Table, startTimestamp, endTimestamp)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("store client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadSpanDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s'`, r.options.Database, r.options.Table, startTimestamp, endTimestamp)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("store client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadDistinctServiceNames(ctx context.Context, dest interface{}) error {
	query := fmt.Sprintf(`SELECT DISTINCT ServiceName as serviceName FROM %s.%s WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		// TODO: log/trace?
		log.Errorf("store client failed to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func NewRepo(opts ...repos.RepoOption) repos.Repo {
	options := repos.NewRepoOptions(opts...)

	r := &sqlRepo{
		options: options,
	}

	return r
}
