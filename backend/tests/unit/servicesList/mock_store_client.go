package serviceslist

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

type mockStoreClient struct {
	options        store.ClientOptions
	readImpl       func() error
	xs             []interface{}
	readCalledWith []map[string]interface{}
	mtx            sync.RWMutex
}

func (c *mockStoreClient) Options() store.ClientOptions {
	return c.options
}

func (c *mockStoreClient) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	c.mtx.Lock()

	args := map[string]interface{}{
		"str": str,
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

func (c *mockStoreClient) ReadCalledWith() []map[string]interface{} {
	return c.readCalledWith
}

func (c *mockStoreClient) ResetCalledWith() {
	c.readCalledWith = []map[string]interface{}{}
}

func NewClient(readImpl func() error, xs ...interface{}) store.Client {
	c := &mockStoreClient{
		options:        store.NewClientOptions(),
		readImpl:       readImpl,
		xs:             xs,
		readCalledWith: []map[string]interface{}{},
		mtx:            sync.RWMutex{},
	}

	return c
}
