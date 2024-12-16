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
	// ReadServices(ctx context.Context, dest interface{}, ...) error
	// ReadServiceMapDependencies(ctx context.Context, dest interface{}, ...) error
	ReadServicesList(ctx context.Context, dest interface{}) error
}
