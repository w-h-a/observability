package memory

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

type client struct {
	options store.ClientOptions
}

func (c *client) Options() store.ClientOptions {
	return c.options
}

func (c *client) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	// TODO
	return nil
}

func NewClient(opts ...store.ClientOption) store.Client {
	options := store.NewClientOptions(opts...)

	return &client{options}
}
