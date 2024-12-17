package unit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
	"github.com/w-h-a/trace-blame/backend/src/clients/store/mock"
)

type testCase struct {
	when           string
	client         store.Client
	then           string
	readCalledWith []map[string]interface{}
	payload        string
}

func TestServicesList(t *testing.T) {
	testCases := GenerateTestCases()

	for _, testCase := range testCases {
		t.Run(testCase.when, func(t *testing.T) {
			var bs []byte
			var err error

			httpServer := src.AppFactory(testCase.client)

			mockStoreClient := testCase.client.(*mock.MockStoreClient)

			err = httpServer.Run()
			require.NoError(t, err)

			bs, err = httputils.HttpGet(fmt.Sprintf("%s/api/v1/services/list", httpServer.Options().Address))
			require.NoError(t, err)

			err = httpServer.Stop()
			require.NoError(t, err)

			t.Log(testCase.then)

			calls := mockStoreClient.ReadCalledWith()

			require.Equal(t, len(calls), 1)

			for i, call := range calls {
				require.Equal(t, testCase.readCalledWith[i], call)
			}

			require.Equal(t, testCase.payload, string(bs))
		})
	}

}

func GenerateTestCases() []testCase {
	successClient := mock.NewClient(
		mock.ClientWithReadImpl(func() error { return nil }),
		mock.ClientWithData("redis", "postgres", "customer", "payments"),
	)

	failureClient := mock.NewClient(
		mock.ClientWithReadImpl(func() error { return fmt.Errorf("failed to process sql query") }),
	)

	testCases := []testCase{
		{
			when:   "when: we get a request to list services and the store client makes a successful call",
			client: successClient,
			then:   "then: we send back the slice of service names",
			readCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT DISTINCT ServiceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
					"additional": []interface{}(nil),
				},
			},
			payload: `["redis","postgres","customer","payments"]`,
		},
		{
			when:   "when: we get a request to list services and the store client fails to make the call",
			client: failureClient,
			then:   "then: we send back an internal server error message",
			readCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT DISTINCT ServiceName FROM . WHERE toDate(Timestamp) > now() - INTERVAL 1 DAY`,
					"additional": []interface{}(nil),
				},
			},
			payload: `{"id":"Services.GetServicesList","code":500,"detail":"failed to retrieve services list: failed to process query","status":"Internal Server Error"}`,
		},
	}

	return testCases
}
