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
	ReadServerSpansOfServices(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadErrServerSpansOfServices(ctx context.Context, dest interface{}, startTimestamp, endTimestamp string) error
	ReadDistinctServiceNames(ctx context.Context, dest interface{}) error
}
