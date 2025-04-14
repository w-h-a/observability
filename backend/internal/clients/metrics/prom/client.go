package prom

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/w-h-a/observability/backend/internal/clients/metrics"
	"github.com/w-h-a/pkg/telemetry/log"
)

type client struct {
	options metrics.ClientOptions
	v1.API
}

func (c *client) Options() metrics.ClientOptions {
	return c.options
}

func (c *client) Read(ctx context.Context, dest interface{}, str string, start time.Time, end time.Time, step time.Duration) error {
	sliceType := reflect.TypeOf(dest).Elem()

	ptr := reflect.ValueOf(dest)

	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("dest is not a pointer or reference")
	}

	if ptr.IsNil() {
		return fmt.Errorf("dest is nil")
	}

	slice := reflect.Indirect(ptr)

	slice.SetLen(0)

	rangeQuery := v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}

	result, _, err := c.QueryRange(ctx, str, rangeQuery)
	if err != nil {
		return err
	}

	matrix, ok := result.(model.Matrix)
	if !ok {
		return fmt.Errorf("unexpected model type: %s", result.Type())
	}

	for _, stream := range matrix {
		for _, pair := range stream.Values {
			typ := sliceType.Elem()
			val := reflect.New(typ).Elem()
			val.FieldByName("Timestamp").SetInt(int64(pair.Timestamp))
			val.FieldByName("Value").SetFloat(float64(pair.Value))
			slice.Set(reflect.Append(slice, val))
		}
	}

	return nil
}

func NewClient(opts ...metrics.ClientOption) metrics.Client {
	options := metrics.NewClientOptions(opts...)

	source := options.Addrs[0]
	if _, err := url.Parse(source); err != nil {
		log.Fatal(err)
	}

	c, err := api.NewClient(api.Config{
		Address: source,
	})
	if err != nil {
		log.Fatal(err)
	}

	a := v1.NewAPI(c)

	return &client{options, a}
}
