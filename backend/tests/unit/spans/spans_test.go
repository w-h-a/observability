package spans

import (
	"fmt"
	"testing"

	"github.com/w-h-a/trace-blame/backend/src/clients/traces/mock"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
	"github.com/w-h-a/trace-blame/backend/tests/unit"
)

func TestSpans(t *testing.T) {
	tagQueryKey := "http.method"
	tagQueryValue := "GET"

	successClient := mock.NewClient(
		mock.RepoClientWithReadImpl(func() error { return nil }),
		mock.RepoClientWithData([][]interface{}{
			{
				&reader.Span{
					Timestamp:    "2024-12-30T00:23:31.474736508Z",
					SpanId:       "4fe6ba8ea64ad354",
					ParentSpanId: "",
					TraceId:      "e5c74cf2d095ad3c85326c90d09312fb",
					ServiceName:  "frontend",
					Name:         "HTTP GET",
					Kind:         "Client",
					StatusCode:   "Ok",
					Duration:     67871375,
					Tags: [][]interface{}{
						{
							"http.response_content_length",
							"59",
						},
						{
							tagQueryKey,
							tagQueryValue,
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
					TraceId:      "e5c74cf2d095ad3c85326c90d09312fb",
					ServiceName:  "route",
					Name:         "/route",
					Kind:         "Server",
					StatusCode:   "Ok",
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
							tagQueryKey,
							tagQueryValue,
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
			When:            "when: we get a request for spans with a tag query that is not allowed",
			Endpoint:        "/api/v1/spans",
			Query:           fmt.Sprintf(`?start=1735770900&end=1735770993&tags=[{"key":"%s","value":"%s","operator":"equal"}]`, tagQueryKey, tagQueryValue),
			Client:          successClient,
			Then:            "then: we return a 400 error response",
			ReadCalledTimes: 0,
			ReadCalledWith:  []map[string]interface{}{},
			Payload:         `{"id":"Spans.GetSpans","code":400,"detail":"failed to parse request: tag operator equal is not supported in tag query","status":"Bad Request"}`,
		},
		{
			When:            "when: we get a request for spans with an _equals_ tag query and the repo client successfully calls the db",
			Endpoint:        "/api/v1/spans",
			Query:           fmt.Sprintf(`?start=1735770900&end=1735770993&tags=[{"key":"%s","value":"%s","operator":"equals"}]`, tagQueryKey, tagQueryValue),
			Client:          successClient,
			Then:            "then: we return the spans",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, StatusCode as statusCode, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM . WHERE timestamp>=? AND timestamp<=? AND SpanAttributes[?]=? ORDER BY timestamp DESC LIMIT 100 OFFSET 0",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000", tagQueryKey, tagQueryValue},
				},
			},
			Payload: `[{"columns":["Time","SpanId","ParentSpanId","TraceId","ServiceName","Name","Kind","StatusCode","Duration","Tags"],"events":[[1735518211474,"4fe6ba8ea64ad354","","e5c74cf2d095ad3c85326c90d09312fb","frontend","HTTP GET","Client","Ok","67871375",[["http.response_content_length","59"],["http.method","GET"],["http.url","http://0.0.0.0:8083/route"],["net.peer.name","0.0.0.0"],["net.peer.port","8083"]]],[1735518211475,"48669d74db753cc2","4fe6ba8ea64ad354","e5c74cf2d095ad3c85326c90d09312fb","route","/route","Server","Ok","67008791",[["http.response_content_length","59"],["net.protocol.version","1.1"],["http.route","/route"],["http.method","GET"],["net.host.name","0.0.0.0"],["net.sock.peer.port","54274"]]]]}]`,
		},
		{
			When:            "when: we get a request for spans with a _contains_ tag query and the repo client successfully calls the db",
			Endpoint:        "/api/v1/spans",
			Query:           fmt.Sprintf(`?start=1735770900&end=1735770993&tags=[{"key":"%s","value":"%s","operator":"contains"}]`, tagQueryKey, tagQueryValue),
			Client:          successClient,
			Then:            "then: we return the spans",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, StatusCode as statusCode, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM . WHERE timestamp>=? AND timestamp<=? AND SpanAttributes[?] ILIKE ? ORDER BY timestamp DESC LIMIT 100 OFFSET 0",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000", tagQueryKey, fmt.Sprintf("%%%s%%", tagQueryValue)},
				},
			},
			Payload: `[{"columns":["Time","SpanId","ParentSpanId","TraceId","ServiceName","Name","Kind","StatusCode","Duration","Tags"],"events":[[1735518211474,"4fe6ba8ea64ad354","","e5c74cf2d095ad3c85326c90d09312fb","frontend","HTTP GET","Client","Ok","67871375",[["http.response_content_length","59"],["http.method","GET"],["http.url","http://0.0.0.0:8083/route"],["net.peer.name","0.0.0.0"],["net.peer.port","8083"]]],[1735518211475,"48669d74db753cc2","4fe6ba8ea64ad354","e5c74cf2d095ad3c85326c90d09312fb","route","/route","Server","Ok","67008791",[["http.response_content_length","59"],["net.protocol.version","1.1"],["http.route","/route"],["http.method","GET"],["net.host.name","0.0.0.0"],["net.sock.peer.port","54274"]]]]}]`,
		},
		{
			When:            "when: we get a request for spans with an _isnotnull_ tag query and the repo client successfully calls the db",
			Endpoint:        "/api/v1/spans",
			Query:           fmt.Sprintf(`?start=1735770900&end=1735770993&tags=[{"key":"%s","value":"%s","operator":"isnotnull"}]`, tagQueryKey, tagQueryValue),
			Client:          successClient,
			Then:            "then: we return the spans",
			ReadCalledTimes: 1,
			ReadCalledWith: []map[string]interface{}{
				{
					"str":        "SELECT Timestamp as timestamp, SpanId as spanId, ParentSpanId as parentSpanId, TraceId as traceId, ServiceName as serviceName, SpanName as name, SpanKind as kind, StatusCode as statusCode, Duration as duration, arrayMap(key -> tuple(key, SpanAttributes[key]), SpanAttributes.keys) as tags FROM . WHERE timestamp>=? AND timestamp<=? AND mapContains(SpanAttributes, ?) ORDER BY timestamp DESC LIMIT 100 OFFSET 0",
					"additional": []interface{}{"1735770900000000000", "1735770993000000000", tagQueryKey},
				},
			},
			Payload: `[{"columns":["Time","SpanId","ParentSpanId","TraceId","ServiceName","Name","Kind","StatusCode","Duration","Tags"],"events":[[1735518211474,"4fe6ba8ea64ad354","","e5c74cf2d095ad3c85326c90d09312fb","frontend","HTTP GET","Client","Ok","67871375",[["http.response_content_length","59"],["http.method","GET"],["http.url","http://0.0.0.0:8083/route"],["net.peer.name","0.0.0.0"],["net.peer.port","8083"]]],[1735518211475,"48669d74db753cc2","4fe6ba8ea64ad354","e5c74cf2d095ad3c85326c90d09312fb","route","/route","Server","Ok","67008791",[["http.response_content_length","59"],["net.protocol.version","1.1"],["http.route","/route"],["http.method","GET"],["net.host.name","0.0.0.0"],["net.sock.peer.port","54274"]]]]}]`,
		},
	}

	unit.RunTestCases(t, testCases)
}
