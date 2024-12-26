package servicedependencies

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src"
	"github.com/w-h-a/trace-blame/backend/src/clients/store"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

func TestServiceDependencies(t *testing.T) {
	successClient1 := NewClient(
		func() error { return nil },
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
	)

	successClient2 := NewClient(
		func() error { return nil },
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
			when:            "when: we get a request to retrieve service dependencies and the client makes a successful call to the db",
			endpoint:        "/api/v1/services/dependencies",
			query:           "?start=1734898000&end=1734913905",
			client:          successClient1,
			then:            "then: we send back the service dependencies",
			readCalledTimes: 1,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000'`,
				},
			},
			payload: `[{"parent":"frontend","child":"frontend","callCount":2},{"parent":"frontend","child":"customer","callCount":1}]`,
		},
		{
			when:            "when: we get a request to retrieve service dependencies, the client makes a successful call to the db, but there are no dependencies",
			endpoint:        "/api/v1/services/dependencies",
			query:           "?start=1734898000&end=1734913905",
			client:          successClient2,
			then:            "then: we send back an empty slice",
			readCalledTimes: 1,
			readCalledWith: []map[string]string{
				{
					"str": `SELECT SpanId as spanId, ParentSpanId as parentSpanId, ServiceName as serviceName FROM . WHERE Timestamp>='1734898000000000000' AND Timestamp<='1734913905000000000'`,
				},
			},
			payload: `[]`,
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
