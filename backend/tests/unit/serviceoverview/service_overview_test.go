package serviceoverview

import (
	"testing"

	metricsmock "github.com/w-h-a/observability/backend/internal/clients/metrics/mock"
	tracesmock "github.com/w-h-a/observability/backend/internal/clients/traces/mock"
	"github.com/w-h-a/observability/backend/internal/services/reader"
	"github.com/w-h-a/observability/backend/tests/unit"
)

func TestServiceOverview(t *testing.T) {
	successClient1 := tracesmock.NewClient(
		tracesmock.RepoClientWithReadImpl(func() error { return nil }),
		tracesmock.RepoClientWithData([][]interface{}{
			{
				&reader.ServiceOverview{
					Time:         "2024-12-26T16:00:00Z",
					Percentile50: 55396310,
					Percentile95: 65587084,
					Percentile99: 69049850,
					NumCalls:     10,
				},
			},
			{},
		}),
	)

	successClient2 := tracesmock.NewClient(
		tracesmock.RepoClientWithReadImpl(func() error { return nil }),
		tracesmock.RepoClientWithData([][]interface{}{
			{
				&reader.ServiceOverview{
					Time:         "2024-12-26T16:00:00Z",
					Percentile50: 55396310,
					Percentile95: 65587084,
					Percentile99: 69049850,
					NumCalls:     16,
				},
			},
			{
				&reader.ServiceOverview{
					Time:      "2024-12-26T16:00:00Z",
					NumErrors: 4,
				},
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve a service overview and the client makes a successful call to the db but there are no error spans",
			Endpoint:        "/api/v1/service/overview",
			Query:           "?start=1734898000&end=1734913905&service=route&step=60",
			TracesClient:    successClient1,
			MetricsClient:   metricsmock.NewClient(),
			Then:            "then: we send back a slice with the overview for the service without errors",
			ReadCalledTimes: 2,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? GROUP BY time ORDER BY time DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000", "route"},
				},
				{
					"str":        `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as numErrors FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? AND StatusCode='Error' GROUP BY time ORDER BY time DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000", "route"},
				},
			},
			Payload: `[{"timestamp":1735228800000000000,"p50":55396310,"p95":65587084,"p99":69049850,"numCalls":10,"callRate":0.16666667,"numErrors":0,"errorRate":0}]`,
		},
		{
			When:            "when: we get a request to retrieve a service overview and the client makes a successful call to the db and there are error spans",
			Endpoint:        "/api/v1/service/overview",
			Query:           "?start=1734898000&end=1734913905&service=route&step=60",
			TracesClient:    successClient2,
			MetricsClient:   metricsmock.NewClient(),
			Then:            "then: we send back a slice with the overview for the service with errors",
			ReadCalledTimes: 2,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? GROUP BY time ORDER BY time DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000", "route"},
				},
				{
					"str":        `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as numErrors FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND ServiceName=? AND StatusCode='Error' GROUP BY time ORDER BY time DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000", "route"},
				},
			},
			Payload: `[{"timestamp":1735228800000000000,"p50":55396310,"p95":65587084,"p99":69049850,"numCalls":16,"callRate":0.26666668,"numErrors":4,"errorRate":0.06666667}]`,
		},
		{
			When:            "when: we get a request to retrieve a service overview but no service name was queried",
			Endpoint:        "/api/v1/service/overview",
			Query:           "?start=1734898000&end=1734913905&&step=60",
			TracesClient:    successClient2,
			Then:            "then: we send back a 400 error message",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Service.GetServiceOverview","code":400,"detail":"failed to parse request: service param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve a service overview but a step of less than 60 was queried",
			Endpoint:        "/api/v1/service/overview",
			Query:           "?start=1734898000&end=1734913905&service=route&step=10",
			TracesClient:    successClient2,
			Then:            "then: we send back a 400 error message",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Service.GetServiceOverview","code":400,"detail":"failed to parse request: step param is less than 60","status":"Bad Request"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
