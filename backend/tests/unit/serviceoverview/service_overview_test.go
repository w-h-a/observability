package serviceoverview

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func TestServiceOverview(t *testing.T) {
	successClient1 := NewClient(
		func() error { return nil },
		[]interface{}{
			&reader.ServiceOverview{
				Time:         "2024-12-26T16:00:00Z",
				Percentile50: 55396310,
				Percentile95: 65587084,
				Percentile99: 69049850,
				NumCalls:     10,
			},
		},
		[]interface{}{},
	)

	successClient2 := NewClient(
		func() error { return nil },
		[]interface{}{
			&reader.ServiceOverview{
				Time:         "2024-12-26T16:00:00Z",
				Percentile50: 55396310,
				Percentile95: 65587084,
				Percentile99: 69049850,
				NumCalls:     16,
			},
		},
		[]interface{}{
			&reader.ServiceOverview{
				Time:      "2024-12-26T16:00:00Z",
				NumErrors: 4,
			},
		},
	)

	testCases := []struct {
		when            string
		endpoint        string
		query           string
		client          store.Client
		then            string
		readCalledTimes int
		readCalledWith  []map[string]string
		payload         string
	}{
		{
			when:            "when: we get a request to retrieve a service overview and the client makes a successful call to the db but there are no error spans",
			endpoint:        "/api/v1/services/overview",
			query:           "?start=1734898000&end=1734913905&service=route&step=60",
			client:          successClient1,
			then:            "then: we send back a slice with the overview for the service without errors",
			readCalledTimes: 2,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' AND ServiceName='route' GROUP BY time ORDER BY time DESC`,
				},
				{
					"str": `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as numErrors FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' AND ServiceName='route' AND StatusCode='Error' GROUP BY time ORDER BY time DESC`,
				},
			},
			payload: `[{"timestamp":1735228800000000000,"p50":55396310,"p95":65587084,"p99":69049850,"numCalls":10,"callRate":0.16666667,"numErrors":0,"errorRate":0}]`,
		},
		{
			when:            "when: we get a request to retrieve a service overview and the client makes a successful call to the db and there are error spans",
			endpoint:        "/api/v1/services/overview",
			query:           "?start=1734898000&end=1734913905&service=route&step=60",
			client:          successClient2,
			then:            "then: we send back a slice with the overview for the service with errors",
			readCalledTimes: 2,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, quantile(0.99)(Duration) as p99, quantile(0.95)(Duration) as p95, quantile(0.50)(Duration) as p50, count(*) as numCalls FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' AND ServiceName='route' GROUP BY time ORDER BY time DESC`,
				},
				{
					"str": `SELECT toStartOfInterval(Timestamp, INTERVAL 1 minute) as time, count(*) as numErrors FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' AND ServiceName='route' AND StatusCode='Error' GROUP BY time ORDER BY time DESC`,
				},
			},
			payload: `[{"timestamp":1735228800000000000,"p50":55396310,"p95":65587084,"p99":69049850,"numCalls":16,"callRate":0.26666668,"numErrors":4,"errorRate":0.06666667}]`,
		},
		{
			when:            "when: we get a request to retrieve a service overview but no service name was queried",
			endpoint:        "/api/v1/services/overview",
			query:           "?start=1734898000&end=1734913905&&step=60",
			client:          successClient2,
			then:            "then: we send back a 400 error message",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetServiceOverview","code":400,"detail":"failed to parse request: service param missing in query","status":"Bad Request"}`,
		},
		{
			when:            "when: we get a request to retrieve a service overview but a step of less than 60 was queried",
			endpoint:        "/api/v1/services/overview",
			query:           "?start=1734898000&end=1734913905&service=route&step=10",
			client:          successClient2,
			then:            "then: we send back a 400 error message",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetServiceOverview","code":400,"detail":"failed to parse request: step param is less than 60","status":"Bad Request"}`,
		},
	}

	for _, testCase := range testCases {
		var bs []byte
		var err error

		httpServer := src.AppFactory(testCase.client)

		mockStoreClient := testCase.client.(*mockStoreClient)

		t.Run(testCase.when, func(t *testing.T) {
			err = httpServer.Run()
			require.NoError(t, err)

			bs, err = httputils.HttpGet(fmt.Sprintf("%s%s%s", httpServer.Options().Address, testCase.endpoint, testCase.query))
			require.NoError(t, err)

			t.Log(testCase.then)

			calls := mockStoreClient.ReadCalledWith()

			require.Equal(t, testCase.readCalledTimes, len(calls))

			for i, call := range calls {
				require.Equal(t, testCase.readCalledWith[i], call)
			}

			require.Equal(t, testCase.payload, string(bs))

			t.Cleanup(func() {
				err = httpServer.Stop()
				require.NoError(t, err)

				mockStoreClient.ResetCalledWith()
			})
		})
	}
}
