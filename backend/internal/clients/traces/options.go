package traces

import "context"

type RepoOption func(o *RepoOptions)

type RepoOptions struct {
	Client   Client
	Database string
	Table    string
	Context  context.Context
}

func RepoWithClient(c Client) RepoOption {
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

type ClientOption func(o *ClientOptions)

type ClientOptions struct {
	Driver  string
	Addrs   []string
	Context context.Context
}

func ClientWithDriver(driver string) ClientOption {
	return func(o *ClientOptions) {
		o.Driver = driver
	}
}

func ClientWithAddrs(addrs ...string) ClientOption {
	return func(o *ClientOptions) {
		o.Addrs = addrs
	}
}

func NewClientOptions(opts ...ClientOption) ClientOptions {
	options := ClientOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
