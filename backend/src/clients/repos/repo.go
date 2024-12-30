package repos

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
	ReadServiceSpecificOperations(ctx context.Context, dest interface{}, serviceName string) error
	ReadServiceSpecificEndpoints(ctx context.Context, dest interface{}, serviceName string, startTimestamp, endTimestamps string) error
	ReadServiceSpecificServerCalls(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
	ReadServiceSpecificServerErrors(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
	ReadSpanDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadTraceSpecificSpans(ctx context.Context, dest interface{}, traceId string) error
}
