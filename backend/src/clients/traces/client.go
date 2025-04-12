package traces

import "context"

type Client interface {
	Options() ClientOptions
	Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error
}
