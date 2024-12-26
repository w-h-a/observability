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
	ReadDistinctServiceNames(ctx context.Context, dest interface{}) error
	ReadSpanDependencies(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadServiceSpecificServerCalls(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
	ReadServiceSpecificServerErrors(ctx context.Context, dest interface{}, serviceName, interval, startTimestamp, endTimestamp string) error
}
