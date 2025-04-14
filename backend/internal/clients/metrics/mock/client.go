package mock

import (
	"context"
	"time"

	"github.com/w-h-a/observability/backend/internal/clients/metrics"
)

type MockMetricsClient struct {
	options metrics.ClientOptions
}

func (c *MockMetricsClient) Options() metrics.ClientOptions {
	return c.options
}

func (c *MockMetricsClient) Read(ctx context.Context, dest interface{}, str string, start time.Time, end time.Time, step time.Duration) error {
	// TODO

	return nil
}

func NewClient(opts ...metrics.ClientOption) metrics.Client {
	options := metrics.NewClientOptions(opts...)

	c := &MockMetricsClient{
		options: options,
	}

	return c
}
