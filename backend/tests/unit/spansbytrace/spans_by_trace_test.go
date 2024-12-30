package spansbytrace

import (
	"fmt"
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/repos/mock"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestSpansByTrace(t *testing.T) {
	traceId := "e5c74cf2d095ad3c85326c90d09312fb"

	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.Span{
					Timestamp:    "2024-12-30T00:23:31.474736508Z",
					SpanId:       "4fe6ba8ea64ad354",
					ParentSpanId: "",
					TraceId:      traceId,
					ServiceName:  "frontend",
					Name:         "HTTP GET",
					Kind:         "Client",
					Duration:     67871375,
					Tags: [][]interface{}{
						{
							"http.response_content_length",
							"59",
						},
						{
							"http.method",
							"GET",
						},
						{
							"http.url",
							"http://0.0.0.0:8083/route",
						},
						{
							"net.peer.name",
							"0.0.0.0",
						},
						{
							"net.peer.port",
							"8083",
						},
					},
				},
				&reader.Span{
					Timestamp:    "2024-12-30T00:23:31.475082716Z",
					SpanId:       "48669d74db753cc2",
					ParentSpanId: "4fe6ba8ea64ad354",
					TraceId:      traceId,
					ServiceName:  "route",
					Name:         "/route",
					Kind:         "Server",
					Duration:     67008791,
					Tags: [][]interface{}{
						{
							"http.response_content_length",
							"59",
						},
						{
							"net.protocol.version",
							"1.1",
						},
						{
							"http.route",
							"/route",
						},
						{
							"http.method",
							"GET",
						},
						{
							"net.host.name",
							"0.0.0.0",
						},
						{
							"net.sock.peer.port",
							"54274",
						},
					},
				},
			},
		}),
	)

	testCases := []unit.TestCase{
		{
			When:            "when: we get a request to retrieve all spans associated with a traceId but there is no traceId",
			Endpoint:        "/api/v1/spans/trace",
			Query:           "",
			Client:          successClient,
			Then:            "then: we send back a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Spans.GetSpansByTraceId","code":400,"detail":"failed to parse request: traceId param missing in query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request to retrieve all spans associated with a traceId, a traceId is queried, and the repo client successfully fetches the data",
			Endpoint:        "/api/v1/spans/trace",
			Query:           fmt.Sprintf("?traceId=%s", traceId),
			Client:          successClient,
			Then:            "then: we send back the span matrix",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        `SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM . WHERE traceId=?`,
					"additional": []interface{}{traceId},
				},
			},
			Payload: `[{"columns":["Time","SpanId","ParentSpanId","TraceId","ServiceName","Name","Kind","Duration","Tags"],"events":[[1735518211474,"4fe6ba8ea64ad354","","e5c74cf2d095ad3c85326c90d09312fb","frontend","HTTP GET","Client","67871375",[["http.response_content_length","59"],["http.method","GET"],["http.url","http://0.0.0.0:8083/route"],["net.peer.name","0.0.0.0"],["net.peer.port","8083"]]],[1735518211475,"48669d74db753cc2","4fe6ba8ea64ad354","e5c74cf2d095ad3c85326c90d09312fb","route","/route","Server","67008791",[["http.response_content_length","59"],["net.protocol.version","1.1"],["http.route","/route"],["http.method","GET"],["net.host.name","0.0.0.0"],["net.sock.peer.port","54274"]]]]}]`,
		},
	}

	unit.RunTestCases(t, testCases)
}
