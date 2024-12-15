package store

import "context"

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
