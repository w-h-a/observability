package reader

import (
	"strconv"
	"time"

	"github.com/w-h-a/observability/backend/src/clients/traces"
)

// services

type ServicesArgs struct {
	Start     *time.Time
	StartTime string
	End       *time.Time
	EndTime   string
	Period    int
}

type Service struct {
	ServiceName  string  `json:"serviceName" db:"serviceName"`
	Percentile99 float32 `json:"p99" db:"p99"`
	AvgDuration  float32 `json:"avgDuration" db:"avgDuration"`
	NumCalls     int     `json:"numCalls" db:"numCalls"`
	CallRate     float32 `json:"callRate" db:"callRate"`
	NumErrors    int     `json:"numErrors" db:"numErrors"`
	ErrorRate    float32 `json:"errorRate" db:"errorRate"`
}

type ServiceSpanDependency struct {
	SpanId       string `db:"spanId,omitempty"`
	ParentSpanId string `db:"parentSpanId,omitempty"`
	ServiceName  string `db:"serviceName,omitempty"`
}

type ServiceDependency struct {
	Parent    string `json:"parent,omitempty"`
	Child     string `json:"child,omitempty"`
	CallCount int    `json:"callCount,omitempty"`
}

// tags

type TagsArgs struct {
	ServiceName string
}

// operations

type OperationsArgs struct {
	ServiceName string
}

type EndpointsArgs struct {
	ServiceName string
	Start       *time.Time
	StartTime   string
	End         *time.Time
	EndTime     string
}

type Endpoint struct {
	Name         string  `json:"name" db:"name"`
	Percentile50 float32 `json:"p50" db:"p50"`
	Percentile95 float32 `json:"p95" db:"p95"`
	Percentile99 float32 `json:"p99" db:"p99"`
	NumCalls     int     `json:"numCalls" db:"numCalls"`
}

// overview

type OverviewArgs struct {
	ServiceName string
	Start       *time.Time
	StartTime   string
	End         *time.Time
	EndTime     string
	StepSeconds int
}

// service overview

type ServiceOverview struct {
	Timestamp    int64   `json:"timestamp" db:"timestamp"`
	Time         string  `json:"time,omitempty" db:"time,omitempty"`
	Percentile50 float32 `json:"p50" db:"p50"`
	Percentile95 float32 `json:"p95" db:"p95"`
	Percentile99 float32 `json:"p99" db:"p99"`
	NumCalls     int     `json:"numCalls" db:"numCalls"`
	CallRate     float32 `json:"callRate" db:"callRate"`
	NumErrors    int     `json:"numErrors" db:"numErrors"`
	ErrorRate    float32 `json:"errorRate" db:"errorRate"`
}

// http overview

type HttpOverview struct {
	Timestamp   int64   `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Time        string  `json:"time,omitempty" db:"time,omitempty"`
	HttpUrl     string  `json:"httpUrl,omitempty" db:"httpUrl,omitempty"`
	HttpMethod  string  `json:"httpMethod,omitempty" db:"httpMethod,omitempty"`
	AvgDuration float32 `json:"avgDuration,omitempty" db:"avgDuration,omitempty"`
	NumCalls    int     `json:"numCalls,omitempty" db:"numCalls,omitempty"`
	CallRate    float32 `json:"callRate,omitempty" db:"callRate,omitempty"`
	NumErrors   int     `json:"numErrors" db:"numErrors"`
	ErrorRate   float32 `json:"errorRate" db:"errorRate"`
}

// rpc overview

type RpcOverview struct {
	Timestamp   int64   `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Time        string  `json:"time,omitempty" db:"time,omitempty"`
	RpcService  string  `json:"rpcService,omitempty" db:"rpcService,omitempty"`
	RpcMethod   string  `json:"rpcMethod,omitempty" db:"rpcMethod,omitempty"`
	AvgDuration float32 `json:"avgDuration,omitempty" db:"avgDuration,omitempty"`
	NumCalls    int     `json:"numCalls,omitempty" db:"numCalls,omitempty"`
	CallRate    float32 `json:"callRate,omitempty" db:"callRate,omitempty"`
	NumErrors   int     `json:"numErrors" db:"numErrors"`
	ErrorRate   float32 `json:"errorRate" db:"errorRate"`
}

// db overview

type DBOverview struct {
	Timestamp   int64   `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Time        string  `json:"time,omitempty" db:"time,omitempty"`
	DBSystem    string  `json:"dbSystem,omitempty" db:"dbSystem,omitempty"`
	DBOperation string  `json:"dbOperation,omitempty" db:"dbOperation,omitempty"`
	AvgDuration float32 `json:"avgDuration,omitempty" db:"avgDuration,omitempty"`
	NumCalls    int     `json:"numCalls,omitempty" db:"numCalls,omitempty"`
	CallRate    float32 `json:"callRate,omitempty" db:"callRate,omitempty"`
	NumErrors   int     `json:"numErrors" db:"numErrors"`
	ErrorRate   float32 `json:"errorRate" db:"errorRate"`
}

// messaging overview

type MessagingOverview struct {
	Timestamp          int64   `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Time               string  `json:"time,omitempty" db:"time,omitempty"`
	MessagingSystem    string  `json:"messagingSystem,omitempty" db:"messagingSystem,omitempty"`
	MessagingOperation string  `json:"messagingOperation,omitempty" db:"messagingOperation,omitempty"`
	AvgDuration        float32 `json:"avgDuration,omitempty" db:"avgDuration,omitempty"`
	NumCalls           int     `json:"numCalls,omitempty" db:"numCalls,omitempty"`
	CallRate           float32 `json:"callRate,omitempty" db:"callRate,omitempty"`
	NumErrors          int     `json:"numErrors" db:"numErrors"`
	ErrorRate          float32 `json:"errorRate" db:"errorRate"`
}

// traces

type TracesArgs struct {
	Start       *time.Time
	StartTime   string
	End         *time.Time
	EndTime     string
	ServiceName string
	TraceId     string
}

// spans

type SpansArgs struct {
	Start       *time.Time
	StartTime   string
	End         *time.Time
	EndTime     string
	ServiceName string
	Name        string
	Kind        string
	MinDuration string
	MaxDuration string
	TagQueries  []traces.TagQuery
	Order       string
	Limit       int64
	Offset      int64
}

type Span struct {
	Timestamp    string          `db:"timestamp"`
	SpanId       string          `db:"spanId"`
	ParentSpanId string          `db:"parentSpanId"`
	TraceId      string          `db:"traceId"`
	ServiceName  string          `db:"serviceName"`
	Name         string          `db:"name"`
	Kind         string          `db:"kind"`
	StatusCode   string          `db:"statusCode"`
	Duration     int64           `db:"duration"`
	Tags         [][]interface{} `db:"tags"`
}

func (s *Span) ToEventValues() []interface{} {
	timeObj, _ := time.Parse(time.RFC3339Nano, s.Timestamp)

	return []interface{}{int64(timeObj.UnixNano() / 1000000), s.SpanId, s.ParentSpanId, s.TraceId, s.ServiceName, s.Name, s.Kind, s.StatusCode, strconv.FormatInt(s.Duration, 10), s.Tags}
}

type SpanMatrix struct {
	Columns []string        `json:"columns"`
	Events  [][]interface{} `json:"events"`
}

// span agg

type AggregatedSpansArgs struct {
	Dimension         string
	AggregationOption string
	StepSeconds       int
	Start             *time.Time
	StartTime         string
	End               *time.Time
	EndTime           string
	ServiceName       string
	Name              string
	Kind              string
	MinDuration       string
	MaxDuration       string
	TagQueries        []traces.TagQuery
}

type AggregatedSpans struct {
	Timestamp int64   `json:"timestamp,omitempty" db:"timestamp"`
	Time      string  `json:"time,omitempty" db:"time"`
	Value     float32 `json:"value,omitempty" db:"value"`
}
