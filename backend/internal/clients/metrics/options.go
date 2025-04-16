package metrics

import "context"

type RepoOption func(o *RepoOptions)

type RepoOptions struct {
	Client  Client
	Context context.Context
}

func RepoWithClient(c Client) RepoOption {
	return func(o *RepoOptions) {
		o.Client = c
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
	Addrs   []string
	Context context.Context
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

type QueryOption func(o *QueryOptions)

type QueryOptions struct {
	ServiceName string
	Context     context.Context
}

func QueryWithServiceName(name string) QueryOption {
	return func(o *QueryOptions) {
		o.ServiceName = name
	}
}

func NewQueryOptions(opts ...QueryOption) QueryOptions {
	options := QueryOptions{
		Context: context.Background(),
	}

	for _, fn := range opts {
		fn(&options)
	}

	return options
}
