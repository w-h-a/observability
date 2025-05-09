package mock

import (
	"context"

	"github.com/w-h-a/observability/backend/internal/clients/traces"
)

type readImplKey struct{}

func RepoClientWithReadImpl(fun func() error) traces.ClientOption {
	return func(o *traces.ClientOptions) {
		o.Context = context.WithValue(o.Context, readImplKey{}, fun)
	}
}

func GetReadImplFromContext(ctx context.Context) (func() error, bool) {
	fun, ok := ctx.Value(readImplKey{}).(func() error)
	return fun, ok
}

type dataKey struct{}

func RepoClientWithData(xs [][]interface{}) traces.ClientOption {
	return func(o *traces.ClientOptions) {
		o.Context = context.WithValue(o.Context, dataKey{}, xs)
	}
}

func GetDataFromContext(ctx context.Context) ([][]interface{}, bool) {
	xs, ok := ctx.Value(dataKey{}).([][]interface{})
	return xs, ok
}
