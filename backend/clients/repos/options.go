package repos

import "context"

type RepoOption func(o *RepoOptions)

type RepoOptions struct {
	Driver   string
	Addrs    []string
	Database string
	Table    string
	Context  context.Context
}

func RepoWithDriver(driver string) RepoOption {
	return func(o *RepoOptions) {
		o.Driver = driver
	}
}

func RepoWithAddrs(addrs ...string) RepoOption {
	return func(o *RepoOptions) {
		o.Addrs = addrs
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
