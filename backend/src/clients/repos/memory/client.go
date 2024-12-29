package memory

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
)

type client struct {
	options repos.ClientOptions
}

func (c *client) Options() repos.ClientOptions {
	return c.options
}

func (c *client) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	// TODO
	return nil
}

func NewClient(opts ...repos.ClientOption) repos.Client {
	options := repos.NewClientOptions(opts...)

	return &client{options}
}
