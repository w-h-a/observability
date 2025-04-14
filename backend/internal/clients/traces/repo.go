package traces

import (
	"context"
	"errors"
)

var (
	ErrProcessingQuery = errors.New("failed to process query")
)

type Repo interface {
	Options() RepoOptions
	ReadServerCalls(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadServerErrors(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadServiceNames(ctx context.Context, dest interface{}) error
	ReadServiceDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadServiceSpecificTags(ctx context.Context, dest interface{}, serviceName string) error
	ReadServiceSpecificOperations(ctx context.Context, dest interface{}, serviceName string) error
	ReadServiceSpecificEndpoints(ctx context.Context, dest interface{}, serviceName string, startTimestamp, endTimestamps string) error
	ReadServiceSpecificServerCalls(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
	ReadServiceSpecificServerErrors(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
	ReadTraces(ctx context.Context, dest interface{}, startTimestamp, endTimestamp, serviceName, traceId string) error
	ReadSpans(ctx context.Context, dest interface{}, startTimestamp, endTimestamp, serviceName, spanName, spanKind, minDuration, maxDuration string, tagQueries ...TagQuery) error
	ReadAggregatedSpans(ctx context.Context, dest interface{}, dimension, aggregationOption, interval, startTimestamp, endTimestamp, serviceName, spanName, spanKind, minDuration, maxDuration string, tagQueries ...TagQuery) error
}
