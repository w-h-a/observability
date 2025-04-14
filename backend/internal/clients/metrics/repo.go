package metrics

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnsupportedQuery = errors.New("unsupported query")
	ErrProcessingQuery  = errors.New("failed to process query")
)

type Repo interface {
	Options() RepoOptions
	ReadMetrics(ctx context.Context, dest interface{}, dimension string, step time.Duration, start, end time.Time, opts ...QueryOption) error
}
