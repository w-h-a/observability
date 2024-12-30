package serviceoperations

import (
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos/mock"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestServiceOperations(t *testing.T) {
	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				"FindDriverIDs", "GetDriver",
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve a service's operations and the client makes a successful call to the db",
			Endpoint:        "/api/v1/service/operations",
			Query:           "?service=driver",
			Client:          successClient,
			Then:            "then: we send back a slice of the operations for the service",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]string{
				{
					"str": `SELECT DISTINCT SpanName as spanName FROM . WHERE ServiceName='driver' AND toDate(Timestamp) > now() - INTERVAL 1 DAY`,
				},
			},
			Payload: `["FindDriverIDs","GetDriver"]`,
		},
		{
			When:            "when: we get a request to retrieve a service's operations without a service param",
			Endpoint:        "/api/v1/service/operations",
			Query:           "?services=driver",
			Client:          successClient,
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]string{},
			Payload:         `{"id":"Service.GetOperations","code":400,"detail":"failed to parse request: service param missing in query","status":"Bad Request"}`,
		},
	}

	unit.RunTestCases(t, testCases)
}
