package sql

import (
	"context"
	"fmt"

	"github.com/w-h-a/trace-blame/backend/src/repos"
)

type sqlRepo struct {
	options repos.RepoOptions
}

func (r *sqlRepo) Options() repos.RepoOptions {
	return r.options
}

func (r *sqlRepo) ReadServicesList(ctx context.Context, dest interface{}) error {
	query := fmt.Sprintf(`SELECT DISTINCT ServiceName FROM %s.%s WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`, r.options.Database, r.options.Table)

	if err := r.options.Client.Read(ctx, dest, query); err != nil {
		// TODO: log/trace?
		return fmt.Errorf("failed to process sql query")
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
