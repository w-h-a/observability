package reader

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/clients/repos"
)

type ReaderOption func(o *ReaderOptions)

type ReaderOptions struct {
	Repo    repos.Repo
	Context context.Context
}

func ReaderWithRepo(repo repos.Repo) ReaderOption {
	return func(o *ReaderOptions) {
		o.Repo = repo
	}
}

func NewReaderOptions(opts ...ReaderOption) ReaderOptions {
	options := ReaderOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
