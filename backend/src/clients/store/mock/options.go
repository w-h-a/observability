package mock

import (
	"context"

	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

type readImplKey struct{}
type dataKey struct{}

func ClientWithReadImpl(impl func() error) store.ClientOption {
	return func(o *store.ClientOptions) {
		o.Context = context.WithValue(o.Context, readImplKey{}, impl)
	}
}

func GetReadImplFromContext(ctx context.Context) (func() error, bool) {
	i, ok := ctx.Value(readImplKey{}).(func() error)
	return i, ok
}

func ClientWithData(xs ...interface{}) store.ClientOption {
	return func(o *store.ClientOptions) {
		o.Context = context.WithValue(o.Context, dataKey{}, xs)
	}
}

func GetDataFromContext(ctx context.Context) ([]interface{}, bool) {
	xs, ok := ctx.Value(dataKey{}).([]interface{})
	return xs, ok
}
