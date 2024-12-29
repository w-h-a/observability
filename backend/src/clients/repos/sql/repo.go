package sql

import (
	"context"
	"fmt"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
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
		log.Errorf("repo client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServerErrors(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT ServiceName as serviceName, count(*) as numErrors FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' AND StatusCode='Error' GROUP BY serviceName`, r.options.Database, r.options.Table, startTimestamp, endTimestamp)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadDistinctServiceNames(ctx context.Context, dest interface{}) error {
	query := fmt.Sprintf(`SELECT DISTINCT ServiceName as serviceName FROM %s.%s WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		// TODO: log/trace?
		log.Errorf("repo client failed to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadSpanDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s'`, r.options.Database, r.options.Table, startTimestamp, endTimestamp)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificOperations(ctx context.Context, dest interface{}, serviceName string) error {
	query := fmt.Sprintf(`SELECT DISTINCT SpanName as spanName FROM %s.%s WHERE ServiceName='%s' AND toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table, serviceName)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificEndpoints(ctx context.Context, dest interface{}, serviceName string, startTimestamp, endTimestamps string) error {
	query := fmt.Sprintf(`SELECT quantile(0.5)(Duration) as p50, quantile(0.95)(Duration) as p95, quantile(0.99)(Duration) as p99, count(*) as numCalls, SpanName as name FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' and ServiceName='%s' GROUP BY name`, r.options.Database, r.options.Table, startTimestamp, endTimestamps, serviceName)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client fail to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificServerCalls(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' AND ServiceName='%s' GROUP BY time ORDER BY time DESC`, interval, r.options.Database, r.options.Table, startTimestamp, endTimestamp, serviceName)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return repos.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificServerErrors(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, count(*) as numErrors FROM %s.%s WHERE Timestamp>='%s' AND Timestamp<='%s' AND SpanKind='Server' AND ServiceName='%s' AND StatusCode='Error' GROUP BY time ORDER BY time DESC`, interval, r.options.Database, r.options.Table, startTimestamp, endTimestamp, serviceName)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client failed to read: %v", err)
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
