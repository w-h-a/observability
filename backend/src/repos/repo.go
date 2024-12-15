package repos

import "context"

type Repo interface {
	Options() RepoOptions
	// ReadServices(ctx context.Context, dest interface{}, ...) error
	// ReadServiceMapDependencies(ctx context.Context, dest interface{}, ...) error
	ReadServicesList(ctx context.Context, dest interface{}) error
}
