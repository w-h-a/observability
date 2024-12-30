package services

import (
	"fmt"
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos/mock"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestServices(t *testing.T) {
	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.Service{
					ServiceName:  "route",
					Percentile99: 723817200,
					AvgDuration:  349520830,
					NumCalls:     4000,
				},
				&reader.Service{
					ServiceName:  "customer",
					Percentile99: 299650080,
					AvgDuration:  284536770,
					NumCalls:     40,
				},
			},
			{
				&reader.Service{
					ServiceName: "route",
					NumErrors:   1,
				},
				&reader.Service{
					ServiceName: "customer",
					NumErrors:   0,
				},
			},
		}),
	)

	failureClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return fmt.Errorf("failed to process sql query") }),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request without start or end query params",
			Endpoint:        "/api/v1/services",
			Query:           "",
			Client:          successClient,
			Then:            "then: we send back a 400 error message",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Services.GetServices","code":400,"detail":"failed to parse request: start param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request without end query params",
			Endpoint:        "/api/v1/services",
			Query:           "?start=1734898000",
			Client:          failureClient,
			Then:            "then: we send back a 400 error message",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Services.GetServices","code":400,"detail":"failed to parse request: end param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve service data and the client makes a successful call to the db",
			Endpoint:        "/api/v1/services",
			Query:           "?start=1734898000&end=1734913905",
			Client:          successClient,
			Then:            "then: we send back the slice of services",
			ReadCalledTimes: 2,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000"},
				},
				{
					"str":        `SELECT ServiceName as serviceName, count(*) as numErrors FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' AND StatusCode='Error' GROUP BY serviceName`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000"},
				},
			},
			Payload: `[{"serviceName":"route","p99":723817200,"avgDuration":349520830,"numCalls":4000,"callRate":0.25149325,"numErrors":1,"errorRate":0.00006287331},{"serviceName":"customer","p99":299650080,"avgDuration":284536770,"numCalls":40,"callRate":0.0025149323,"numErrors":0,"errorRate":0}]`,
		},
		{
			When:            "when: we get a request to retrieve service data and the client fails to make the call to the db",
			Endpoint:        "/api/v1/services",
			Query:           "?start=1734898000&end=1734913905",
			Client:          failureClient,
			Then:            "then: we send back a 500 error message",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM . WHERE Timestamp>=? AND Timestamp<=? AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`,
					"additional": []interface{}{"1734898000000000000", "1734913905000000000"},
				},
			},
			Payload: `{"id":"Services.GetServices","code":500,"detail":"failed to retrieve services: failed to process query","status":"Internal Server Error"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
