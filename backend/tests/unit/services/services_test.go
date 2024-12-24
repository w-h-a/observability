package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func TestServices(t *testing.T) {
	successClient := NewClient(
		func() error { return nil },
		[]interface{}{
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
		[]interface{}{
			&reader.Service{
				ServiceName: "route",
				NumErrors:   1,
			},
			&reader.Service{
				ServiceName: "customer",
				NumErrors:   0,
			},
		},
	)

	failureClient := NewClient(
		func() error { return fmt.Errorf("failed to process sql query") },
		[]interface{}{},
		[]interface{}{},
	)

	testCases := []struct {
		when            string
		query           string
		client          store.Client
		then            string
		readCalledTimes int
		readCalledWith  []map[string]string
		payload         string
	}{
		{
			when:            "when: we get a request without start or end query params",
			query:           "",
			client:          successClient,
			then:            "then: we send back a 400 error message",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetServices","code":400,"detail":"failed to parse request: start param missing in query","status":"Bad Request"}`,
		},
		{
			when:            "when: we get a request without end query params",
			query:           "?start=1734898000",
			client:          failureClient,
			then:            "then: we send back a 400 error message",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetServices","code":400,"detail":"failed to parse request: end param missing in query","status":"Bad Request"}`,
		},
		{
			when:            "when: we get a request to retrieve service data and the client makes a successful call to the db",
			query:           "?start=1734898000&end=1734913905",
			client:          successClient,
			then:            "then: we send back the slice of services",
			readCalledTimes: 2,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`,
				},
				{
					"str": `SELECT ServiceName as serviceName, count(*) as numErrors FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' AND StatusCode='Error' GROUP BY serviceName`,
				},
			},
			payload: `[{"serviceName":"route","p99":723817200,"avgDuration":349520830,"numCalls":4000,"callRate":0.25149325,"numErrors":1,"errorRate":0.00006287331},{"serviceName":"customer","p99":299650080,"avgDuration":284536770,"numCalls":40,"callRate":0.0025149323,"numErrors":0,"errorRate":0}]`,
		},
		{
			when:            "when: we get a request to retrieve service data and the client fails to make the call to the db",
			query:           "?start=1734898000&end=1734913905",
			client:          failureClient,
			then:            "then: we send back a 500 error message",
			readCalledTimes: 1,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT ServiceName as serviceName, quantile(0.99)(Duration) as p99, avg(Duration) as avgDuration, count(*) as numCalls FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000' AND SpanKind='Server' GROUP BY serviceName ORDER BY p99 DESC`,
				},
			},
			payload: `{"id":"Services.GetServices","code":500,"detail":"failed to retrieve services: failed to process query","status":"Internal Server Error"}`,
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

			bs, err = httputils.HttpGet(fmt.Sprintf("%s/api/v1/services%s", httpServer.Options().Address, testCase.query))
			require.NoError(t, err)

			t.Log(testCase.then)

			calls := mockStoreClient.ReadCalledWith()

			require.Equal(t, len(calls), testCase.readCalledTimes)

			for i, call := range calls {
				require.Equal(t, testCase.readCalledWith[i], call)
			}

			require.Equal(t, testCase.payload, string(bs))
		})

		t.Cleanup(func() {
			err = httpServer.Stop()
			require.NoError(t, err)

			mockStoreClient.ResetCalledWith()
		})
	}
}
