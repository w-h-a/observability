package serviceendpoints

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/repos"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func TestServiceEndpoints(t *testing.T) {
	successClient := NewClient(
		func() error { return nil },
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
	)

	testCases := []struct {
		when            string
		endpoint        string
		query           string
		client          repos.Client
		then            string
		readCalledTimes int
		readCalledWith  []map[string]string
		payload         string
	}{
		{
			when:            "when: we get a request to retrieve a service's endpoints and the client makes a successful call to the db",
			endpoint:        "/api/v1/services/endpoints",
			query:           "?service=frontend&start=1735254100&end=1735254931",
			client:          successClient,
			then:            "then: we send back a slice of the endpoints for the service",
			readCalledTimes: 1,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT quantile(0.5)(Duration) as p50, quantile(0.95)(Duration) as p95, quantile(0.99)(Duration) as p99, count(*) as numCalls, SpanName as name FROM . WHERE Timestamp>='1735254100000000000' AND Timestamp<='1735254931000000000' AND SpanKind='Server' and ServiceName='frontend' GROUP BY name`,
				},
			},
			payload: `[{"name":"/config","p50":70979,"p95":143176.75,"p99":148701.75,"numCalls":6},{"name":"/dispatch","p50":751639360,"p95":751639360,"p99":751639360,"numCalls":6}]`,
		},
		{
			when:            "when: we get a request to retrieve a service's endpoints without start or end params",
			endpoint:        "/api/v1/services/endpoints",
			query:           "?service=driver",
			client:          successClient,
			then:            "then: we send back a 400 error response",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetEndpoints","code":400,"detail":"failed to parse request: start param missing in query","status":"Bad Request"}`,
		},
	}

	for _, testCase := range testCases {
		var bs []byte
		var err error

		httpServer := src.AppFactory(testCase.client)

		mockStoreClient := testCase.client.(*mockRepoClient)

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
