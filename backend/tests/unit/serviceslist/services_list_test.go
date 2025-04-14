package serviceslist

import (
	"fmt"
	"testing"

	metricsmock "github.com/w-h-a/observability/backend/internal/clients/metrics/mock"
	tracesmock "github.com/w-h-a/observability/backend/internal/clients/traces/mock"
	"github.com/w-h-a/observability/backend/tests/unit"
)

func TestServicesList(t *testing.T) {
	successClient := tracesmock.NewClient(
		tracesmock.RepoClientWithReadImpl(func() error { return nil }),
		tracesmock.RepoClientWithData([][]interface{}{
			{
				"redis", "postgres", "customer", "payments",
			},
		}),
	)

	failureClient := tracesmock.NewClient(
		tracesmock.RepoClientWithReadImpl(func() error { return fmt.Errorf("failed to process sql query") }),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to list services and the store client makes a successful call to the db",
			Endpoint:        "/api/v1/services/list",
			TracesClient:    successClient,
			MetricsClient:   metricsmock.NewClient(),
			Then:            "then: we send back the slice of service names",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT DISTINCT ServiceName as serviceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
					"additional": []interface{}(nil),
				},
			},
			Payload: `["redis","postgres","customer","payments"]`,
		},
		{
			When:            "when: we get a request to list services and the store client fails to make the call to the db",
			Endpoint:        "/api/v1/services/list",
			TracesClient:    failureClient,
			MetricsClient:   metricsmock.NewClient(),
			Then:            "then: we send back an internal server error message",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT DISTINCT ServiceName as serviceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
					"additional": []interface{}(nil),
				},
			},
			Payload: `{"id":"Services.GetServicesList","code":500,"detail":"failed to retrieve services list: failed to process query","status":"Internal Server Error"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
