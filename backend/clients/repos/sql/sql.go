package sql

import (
	"net/url"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"

	"github.com/w-h-a/pkg/telemetry/log"
	"github.com/w-h-a/trace-blame/backend/clients/repos"
)

type sqlRepo struct {
	options repos.RepoOptions
	client  *sqlx.DB
}

func (r *sqlRepo) Options() repos.RepoOptions {
	return r.options
}

func (r *sqlRepo) Read(dest interface{}, str string, additional ...interface{}) error {
	return r.client.Select(dest, str, additional...)
}

func (r *sqlRepo) configure() error {
	source := r.options.Addrs[0]
	if _, err := url.Parse(source); err != nil {
		return err
	}

	client, err := sqlx.Open(r.options.Driver, source)
	if err != nil {
		return err
	}

	if err := client.Ping(); err != nil {
		return err
	}

	r.client = client

	return nil
}

func NewRepo(opts ...repos.RepoOption) repos.Repo {
	options := repos.NewRepoOptions(opts...)

	r := &sqlRepo{
		options: options,
	}

	if err := r.configure(); err != nil {
		log.Fatal(err)
	}

	return r
}
