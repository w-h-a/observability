package metrics

import (
	"context"
	"time"
)

type Client interface {
	Options() ClientOptions
	Read(ctx context.Context, dest interface{}, str string, start time.Time, end time.Time, step time.Duration) error
}
