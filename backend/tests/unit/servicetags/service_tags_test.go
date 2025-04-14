package servicetags

import (
	"testing"

	"github.com/w-h-a/observability/backend/src/clients/traces/mock"
	"github.com/w-h-a/observability/backend/tests/unit"
)

func TestServiceTags(t *testing.T) {
	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				"rpc.service", "rpc.method", "rpc.system", "net.sock.peer.addr", "net.sock.peer.port",
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve a service's tags but no service is queried",
			Endpoint:        "/api/v1/service/tags",
			Client:          successClient,
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Service.GetTags","code":400,"detail":"failed to parse request: service param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve a service's tags, a service is queried, and the repo client successfully makes the call to the db",
			Endpoint:        "/api/v1/service/tags",
			Query:           "?service=driver",
			Client:          successClient,
			Then:            "then: we send back a slice of unique tag key values for this service",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT DISTINCT arrayJoin(SpanAttributes.keys) as tags FROM . WHERE ServiceName=? AND toDate(Timestamp) > now() - INTERVAL 1 DAY`,
					"additional": []interface{}{"driver"},
				},
			},
			Payload: `["rpc.service","rpc.method","rpc.system","net.sock.peer.addr","net.sock.peer.port"]`,
		},
	}

	unit.RunTestCases(t, testCases)
}
