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
}
