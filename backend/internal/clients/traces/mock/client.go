package mock

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/w-h-a/observability/backend/internal/clients/traces"
)

type MockTracesClient struct {
	options         traces.ClientOptions
	readImpl        func() error
	xs              [][]interface{}
	readCalledTimes int
	readCalledWith  []map[string]interface{}
	mtx             sync.RWMutex
}

func (c *MockTracesClient) Options() traces.ClientOptions {
	return c.options
}

func (c *MockTracesClient) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	c.mtx.Lock()

	args := map[string]interface{}{
		"str":        str,
		"additional": additional,
	}

	c.readCalledTimes += 1

	c.readCalledWith = append(c.readCalledWith, args)

	c.mtx.Unlock()

	if err := c.readImpl(); err != nil {
		return err
	}

	ptr := reflect.ValueOf(dest)

	if ptr.Kind() != reflect.Ptr {
		return fmt.Errorf("dest is not a pointer or reference")
	}

	if ptr.IsNil() {
		return fmt.Errorf("dest is nil")
	}

	slice := reflect.Indirect(ptr)

	slice.SetLen(0)

	c.mtx.RLock()

	xs := c.xs[c.readCalledTimes-1]

	c.mtx.RUnlock()

	for _, x := range xs {
		slice.Set(reflect.Append(slice, reflect.ValueOf(x)))
	}

	return nil
}

func (c *MockTracesClient) ReadCalledWith() []map[string]interface{} {
	return c.readCalledWith
}

func (c *MockTracesClient) ResetCalledWith() {
	c.readCalledTimes = 0
	c.readCalledWith = []map[string]interface{}{}
}

func NewClient(opts ...traces.ClientOption) traces.Client {
	options := traces.NewClientOptions(opts...)

	c := &MockTracesClient{
		options:        options,
		readCalledWith: []map[string]interface{}{},
		mtx:            sync.RWMutex{},
	}

	if fun, ok := GetReadImplFromContext(options.Context); ok {
		c.readImpl = fun
	}

	if xs, ok := GetDataFromContext(options.Context); ok {
		c.xs = xs
	}

	return c
}
