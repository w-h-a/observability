package mock

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

type MockStoreClient struct {
	options        store.ClientOptions
	readImpl       func() error
	xs             []interface{}
	readCalledWith []map[string]interface{}
	mtx            sync.RWMutex
}

func (c *MockStoreClient) Options() store.ClientOptions {
	return c.options
}

func (c *MockStoreClient) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	c.mtx.Lock()

	args := map[string]interface{}{
		"str":        str,
		"additional": additional,
	}

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

	for _, x := range c.xs {
		slice.Set(reflect.Append(slice, reflect.ValueOf(x)))
	}

	return nil
}

func (c *MockStoreClient) ReadCalledWith() []map[string]interface{} {
	return c.readCalledWith
}

func NewClient(opts ...store.ClientOption) store.Client {
	options := store.NewClientOptions(opts...)

	c := &MockStoreClient{
		options:        options,
		readCalledWith: []map[string]interface{}{},
		mtx:            sync.RWMutex{},
	}

	if readImpl, ok := GetReadImplFromContext(options.Context); ok {
		c.readImpl = readImpl
	}

	if xs, ok := GetDataFromContext(options.Context); ok {
		c.xs = xs
	}

	return c
}
