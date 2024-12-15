package repos

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

type RepoOption func(o *RepoOptions)

type RepoOptions struct {
	Client   store.Client
	Database string
	Table    string
	Context  context.Context
}

func RepoWithClient(c store.Client) RepoOption {
	return func(o *RepoOptions) {
		o.Client = c
	}
}

func RepoWithDatabase(db string) RepoOption {
	return func(o *RepoOptions) {
		o.Database = db
	}
}

func RepoWithTable(tbl string) RepoOption {
	return func(o *RepoOptions) {
		o.Table = tbl
	}
}

func NewRepoOptions(opts ...RepoOption) RepoOptions {
	options := RepoOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
