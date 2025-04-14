package servicedependencies

import (
	"testing"

	"github.com/w-h-a/observability/backend/src/clients/traces/mock"
	"github.com/w-h-a/observability/backend/src/services/reader"
	"github.com/w-h-a/observability/backend/tests/unit"
)

func TestServiceDependencies(t *testing.T) {
	successClient1 := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.ServiceSpanDependency{
					SpanId:       "da23bbf711742365",
					ParentSpanId: "7443b0b176094d1b",
					ServiceName:  "customer",
				},
				&reader.ServiceSpanDependency{
					SpanId:       "7443b0b176094d1b",
					ParentSpanId: "4eab74f8050def21",
					ServiceName:  "frontend",
				},
				&reader.ServiceSpanDependency{
					SpanId:       "2758719e99111703",
					ParentSpanId: "4eab74f8050def21",
					ServiceName:  "frontend",
				},
				&reader.ServiceSpanDependency{
					SpanId:       "4eab74f8050def21",
					ParentSpanId: "",
					ServiceName:  "frontend",
				},
			},
		}),
	)

	successClient2 := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.ServiceSpanDependency{
					SpanId:       "da23bbf711742365",
					ParentSpanId: "7443b0b176094d1b",
					ServiceName:  "customer",
				},
				&reader.ServiceSpanDependency{
					SpanId:       "4eab74f8050def21",
					ParentSpanId: "",
					ServiceName:  "frontend",
				},
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve service dependencies and the client makes a successful call to the db",
			Endpoint:        "/api/v1/services/dependencies",
			Query:           "?start=1734898000&end=1734913905",
			Client:          successClient1,
			Then:            "then: we send back the service dependencies",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM . WHERE Timestamp>=? AND Timestamp<=?`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000"},
				},
			},
			Payload: `[{"parent":"frontend","child":"frontend","callCount":2},{"parent":"frontend","child":"customer","callCount":1}]`,
		},
		{
			When:            "when: we get a request to retrieve service dependencies, the client makes a successful call to the db, but there are no dependencies",
			Endpoint:        "/api/v1/services/dependencies",
			Query:           "?start=1734898000&end=1734913905",
			Client:          successClient2,
			Then:            "then: we send back an empty slice",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM . WHERE Timestamp>=? AND Timestamp<=?`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000"},
				},
			},
			Payload: `[]`,
		},
	}

	unit.RunTestCases(t, testCases)
}
