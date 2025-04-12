package sql

import (
	"context"
	"net/url"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/src/clients/traces"
)

type client struct {
	options traces.ClientOptions
	*sqlx.DB
}

func (c *client) Options() traces.ClientOptions {
	return c.options
}

func (c *client) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	return c.Select(dest, str, additional...)
}

func NewClient(opts ...traces.ClientOption) traces.Client {
	options := traces.NewClientOptions(opts...)

	source := options.Addrs[0]
	if _, err := url.Parse(source); err != nil {
		log.Fatal(err)
	}

	c, err := sqlx.Open(options.Driver, source)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Ping(); err != nil {
		log.Fatal(err)
	}

	return &client{options, c}
}
