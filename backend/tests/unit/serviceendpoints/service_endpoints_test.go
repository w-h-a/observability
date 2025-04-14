package serviceendpoints

import (
	"testing"

	"github.com/w-h-a/observability/backend/src/clients/traces/mock"
	"github.com/w-h-a/observability/backend/src/services/reader"
	"github.com/w-h-a/observability/backend/tests/unit"
)

func TestServiceEndpoints(t *testing.T) {
	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.Endpoint{
					Name:         "/config",
					Percentile50: 70979,
					Percentile95: 143176.75,
					Percentile99: 148701.75,
					NumCalls:     6,
				},
				&reader.Endpoint{
					Name:         "/dispatch",
					Percentile50: 751639360,
					Percentile95: 751639360,
					Percentile99: 751639360,
					NumCalls:     6,
				},
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve a service's endpoints and the client makes a successful call to the db",
			Endpoint:        "/api/v1/service/endpoints",
			Query:           "?service=frontend&start=1735254100&end=1735254931",
			Client:          successClient,
			Then:            "then: we send back a slice of the endpoints for the service",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT quantile(0.5)(Duration) as p50, quantile(0.95)(Duration) as p95, quantile(0.99)(Duration) as p99, count(*) as numCalls, SpanName as name FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' and ServiceName=? GROUP BY name`,
					"additional": []interface{}{"1735254100000000000", "1735254931000000000", "frontend"},
				},
			},
			Payload: `[{"name":"/config","p50":70979,"p95":143176.75,"p99":148701.75,"numCalls":6},{"name":"/dispatch","p50":751639360,"p95":751639360,"p99":751639360,"numCalls":6}]`,
		},
		{
			When:            "when: we get a request to retrieve a service's endpoints without start or end params",
			Endpoint:        "/api/v1/service/endpoints",
			Query:           "?service=driver",
			Client:          successClient,
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Service.GetEndpoints","code":400,"detail":"failed to parse request: start param missing in query","status":"Bad Request"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
