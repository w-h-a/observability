package mock

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
)

type MockRepoClient struct {
	options         repos.ClientOptions
	readImpl        func() error
	xs              [][]interface{}
	readCalledTimes int
	readCalledWith  []map[string]string
	mtx             sync.RWMutex
}

func (c *MockRepoClient) Options() repos.ClientOptions {
	return c.options
}

func (c *MockRepoClient) Read(ctx context.Context, dest interface{}, str string, additional ...interface{}) error {
	c.mtx.Lock()

	args := map[string]string{
		"str": str,
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

func (c *MockRepoClient) ReadCalledWith() []map[string]string {
	return c.readCalledWith
}

func (c *MockRepoClient) ResetCalledWith() {
	c.readCalledWith = []map[string]string{}
}

func NewClient(opts ...repos.ClientOption) repos.Client {
	options := repos.NewClientOptions(opts...)

	c := &MockRepoClient{
		options:        options,
		readCalledWith: []map[string]string{},
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
