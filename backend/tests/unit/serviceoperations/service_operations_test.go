package serviceoperations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
)

func TestServiceOperations(t *testing.T) {
	successClient := NewClient(
		func() error { return nil },
		"FindDriverIDs", "GetDriver",
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
			when:            "when: we get a request to retrieve a service's operations and the client makes a successful call to the db",
			endpoint:        "/api/v1/services/operations",
			query:           "?service=driver",
			client:          successClient,
			then:            "then: we send back a slice of the operations for the service",
			readCalledTimes: 1,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT DISTINCT SpanName as spanName FROM . WHERE ServiceName='driver' AND toDate(Timestamp) > now() - INTERVAL 1 DAY`,
				},
			},
			payload: `["FindDriverIDs","GetDriver"]`,
		},
		{
			when:            "when: we get a request to retrieve a service's operations without a service param",
			endpoint:        "/api/v1/services/operations",
			query:           "?services=driver",
			client:          successClient,
			then:            "then: we send back a 400 error response",
			readCalledTimes: 0,
			readCalledWith:  []map[string]string{},
			payload:         `{"id":"Services.GetOperations","code":400,"detail":"failed to parse request: service param missing in query","status":"Bad Request"}`,
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
