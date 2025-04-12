package sql

import (
	"context"
	"fmt"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/src/clients/traces"
)

type sqlRepo struct {
	options traces.RepoOptions
}

func (r *sqlRepo) Options() traces.RepoOptions {
	return r.options
}

func (r *sqlRepo) ReadServerCalls(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM %s.%s WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamp); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServerErrors(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT ServiceName as serviceName, count(*) as numErrors FROM %s.%s WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND StatusCode='Error' GROUP BY serviceName`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamp); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceNames(ctx context.Context, dest interface{}) error {
	query := fmt.Sprintf(`SELECT DISTINCT ServiceName as serviceName FROM %s.%s WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM %s.%s WHERE Timestamp>=? AND Timestamp<=?`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamp); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificTags(ctx context.Context, dest interface{}, serviceName string) error {
	query := fmt.Sprintf(`SELECT DISTINCT arrayJoin(SpanAttributes.keys) as tags FROM %s.%s WHERE ServiceName=? AND toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, serviceName); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificOperations(ctx context.Context, dest interface{}, serviceName string) error {
	query := fmt.Sprintf(`SELECT DISTINCT SpanName as spanName FROM %s.%s WHERE ServiceName=? AND toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, serviceName); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificEndpoints(ctx context.Context, dest interface{}, serviceName string, startTimestamp, endTimestamps string) error {
	query := fmt.Sprintf(`SELECT quantile(0.5)(Duration) as p50, quantile(0.95)(Duration) as p95, quantile(0.99)(Duration) as p99, count(*) as numCalls, SpanName as name FROM %s.%s WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' and ServiceName=? GROUP BY name`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamps, serviceName); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificServerCalls(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM %s.%s WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? GROUP BY time ORDER BY time DESC`, interval, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamp, serviceName); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadServiceSpecificServerErrors(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error {
	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, count(*) as numErrors FROM %s.%s WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? AND StatusCode='Error' GROUP BY time ORDER BY time DESC`, interval, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query, startTimestamp, endTimestamp, serviceName); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadTraces(ctx context.Context, dest interface{}, startTimestamp, endTimestamp, serviceName, traceId string) error {
	query := fmt.Sprintf(`SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, StatusCode as statusCode, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM %s.%s WHERE timestamp>=? AND timestamp<=?`, r.options.Database, r.options.Table)

	args := []interface{}{startTimestamp, endTimestamp}

	if len(serviceName) != 0 {
		query += " AND serviceName=?"
		args = append(args, serviceName)
	}

	if len(traceId) != 0 {
		query += " AND traceId=?"
		args = append(args, traceId)
	}

	if err := r.options.Client.Read(ctx, dest, query, args...); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadSpans(
	ctx context.Context,
	dest interface{},
	startTimestamp,
	endTimestamp,
	serviceName,
	spanName,
	spanKind,
	minDuration,
	maxDuration string,
	tagQueries ...traces.TagQuery,
) error {
	query := fmt.Sprintf(`SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, StatusCode as statusCode, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM %s.%s WHERE timestamp>=? AND timestamp<=?`, r.options.Database, r.options.Table)

	args := []interface{}{startTimestamp, endTimestamp}

	var err error

	query, args, err = r.buildUpSpanQuery(query, args, serviceName, spanName, spanKind, minDuration, maxDuration, tagQueries)
	if err != nil {
		return err
	}

	// TODO: make this configurable
	query += " ORDER BY timestamp DESC LIMIT 100 OFFSET 0"

	if err := r.options.Client.Read(ctx, dest, query, args...); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (r *sqlRepo) ReadAggregatedSpans(
	ctx context.Context,
	dest interface{},
	dimension,
	aggregationOption,
	interval,
	startTimestamp,
	endTimestamp,
	serviceName,
	spanName,
	spanKind,
	minDuration,
	maxDuration string,
	tagQueries ...traces.TagQuery,
) error {
	aggregationQuery := ""

	switch dimension {
	case "duration":
		switch aggregationOption {
		case "avg":
			aggregationQuery = "avg(Duration) as value"
		case "p50":
			aggregationQuery = "quantile(0.50)(Duration) as value"
		case "p95":
			aggregationQuery = "quantile(0.95)(Duration) as value"
		case "p99":
			aggregationQuery = "quantile(0.99)(Duration) as value"
		}
	case "calls":
		aggregationQuery = "count(*) as value"
	}

	query := fmt.Sprintf(`SELECT toStartOfInterval(Timestamp, INTERVAL %s minute) as time, %s FROM %s.%s WHERE Timestamp>=? AND Timestamp<=?`, interval, aggregationQuery, r.options.Database, r.options.Table)

	args := []interface{}{startTimestamp, endTimestamp}

	var err error

	query, args, err = r.buildUpSpanQuery(query, args, serviceName, spanName, spanKind, minDuration, maxDuration, tagQueries)
	if err != nil {
		return err
	}

	query += " GROUP BY time ORDER By time"

	if err := r.options.Client.Read(ctx, dest, query, args...); err != nil {
		log.Errorf("repo client failed to read: %v", err)
		return traces.ErrProcessingQuery
	}

	return nil
}

func (*sqlRepo) buildUpSpanQuery(
	query string,
	args []interface{},
	serviceName string,
	spanName string,
	spanKind string,
	minDuration string,
	maxDuration string,
	tagQueries []traces.TagQuery,
) (string, []interface{}, error) {
	if len(serviceName) != 0 {
		query += " AND ServiceName=?"
		args = append(args, serviceName)
	}

	if len(spanName) != 0 {
		query += " AND SpanName=?"
		args = append(args, spanName)
	}

	if len(spanKind) != 0 {
		query += " AND SpanKind=?"
		args = append(args, spanKind)
	}

	if len(minDuration) != 0 {
		query += " AND Duration>=?"
		args = append(args, minDuration)
	}

	if len(maxDuration) != 0 {
		query += " AND Duration<=?"
		args = append(args, maxDuration)
	}

	for _, tagQuery := range tagQueries {
		if tagQuery.Key == "error" && tagQuery.Value == "true" {
			query += " AND (SpanAttributes['error']='true' OR StatusCode='Error')"
			continue
		}

		switch tagQuery.Operator {
		case "equals":
			query += " AND SpanAttributes[?]=?"
			args = append(args, tagQuery.Key, tagQuery.Value)
		case "contains":
			query += " AND SpanAttributes[?] ILIKE ?"
			args = append(args, tagQuery.Key, fmt.Sprintf("%%%s%%", tagQuery.Value))
		case "isnotnull":
			query += " AND mapContains(SpanAttributes, ?)"
			args = append(args, tagQuery.Key)
		}
	}

	return query, args, nil
}

func NewRepo(opts ...traces.RepoOption) traces.Repo {
	options := traces.NewRepoOptions(opts...)

	r := &sqlRepo{
		options: options,
	}

	return r
}
