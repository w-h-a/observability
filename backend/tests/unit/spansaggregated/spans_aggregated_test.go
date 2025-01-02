package spansaggregated

import (
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos/mock"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestSpansAggregated(t *testing.T) {
	callsClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.AggregatedSpans{
					Time:  "2025-01-02T00:11:00Z",
					Value: 41,
				},
			},
		}),
	)

	durationClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.AggregatedSpans{
					Time:  "2025-01-02T00:11:00Z",
					Value: 83866070,
				},
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve aggregated span data but we do not receive a dimension upon which to aggregate",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60",
			Client:          mock.NewClient(),
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Spans.GetAggregatedSpans","code":400,"detail":"failed to parse request: dimension param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve aggregated span data but we do not receive a valid aggregation option for the dimension",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60&dimension=calls&aggregation=p95",
			Client:          mock.NewClient(),
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Spans.GetAggregatedSpans","code":400,"detail":"failed to parse request: aggregation p95 is not supported in query with dimension calls","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve the rate of calls, aggregated over spans",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60&dimension=calls&aggregation=rate_per_sec",
			Client:          callsClient,
			Then:            "then: we send back the aggregated data",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as value FROM . WHERE Timestamp>=? AND Timestamp<=? GROUP BY time ORDER By time",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000"},
				},
			},
			Payload: `[{"timestamp":1735776660000000000,"value":0.68333334}]`,
		},
		{
			When:            "when: we get a request to retrieve the count of calls, aggregated over spans",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60&dimension=calls&aggregation=count",
			Client:          callsClient,
			Then:            "then: we send back the aggregated data",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as value FROM . WHERE Timestamp>=? AND Timestamp<=? GROUP BY time ORDER By time",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000"},
				},
			},
			Payload: `[{"timestamp":1735776660000000000,"value":41}]`,
		},
		{
			When:            "when: we get a request to retrieve the avg of duration, aggregated over spans",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60&dimension=duration&aggregation=avg",
			Client:          durationClient,
			Then:            "then: we send back the aggregated data",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, avg(Duration) as value FROM . WHERE Timestamp>=? AND Timestamp<=? GROUP BY time ORDER By time",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000"},
				},
			},
			Payload: `[{"timestamp":1735776660000000000,"value":83866070}]`,
		},
		{
			When:            "when: we get a request to retrieve the p95 of duration, aggregated over spans",
			Endpoint:        "/api/v1/spans/aggregated",
			Query:           "?start=1735770900&end=1735770993&step=60&dimension=duration&aggregation=p95",
			Client:          durationClient,
			Then:            "then: we send back the aggregated data",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, quantile(0.95)(Duration) as value FROM . WHERE Timestamp>=? AND Timestamp<=? GROUP BY time ORDER By time",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000"},
				},
			},
			Payload: `[{"timestamp":1735776660000000000,"value":83866070}]`,
		},
	}

	unit.RunTestCases(t, testCases)
}
