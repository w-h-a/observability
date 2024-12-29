package serviceslist

import (
	"fmt"
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos/mock"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestServicesList(t *testing.T) {
	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				"redis", "postgres", "customer", "payments",
			},
		}),
	)

	failureClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return fmt.Errorf("failed to process sql query") }),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to list services and the store client makes a successful call to the db",
			Endpoint:        "/api/v1/services/list",
			Client:          successClient,
			Then:            "then: we send back the slice of service names",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]string{
				{
					"str": `SELECT DISTINCT ServiceName as serviceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
				},
			},
			Payload: `["redis","postgres","customer","payments"]`,
		},
		{
			When:            "when: we get a request to list services and the store client fails to make the call to the db",
			Endpoint:        "/api/v1/services/list",
			Client:          failureClient,
			Then:            "then: we send back an internal server error message",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]string{
				{
					"str": `SELECT DISTINCT ServiceName as serviceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
				},
			},
			Payload: `{"id":"Services.GetServicesList","code":500,"detail":"failed to retrieve services list: failed to process query","status":"Internal Server Error"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
