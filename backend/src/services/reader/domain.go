package reader

import "time"

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
	SpanId       string `json:"spanId,omitempty" db:"spanId,omitempty"`
	ParentSpanId string `json:"parentSpanId,omitempty" db:"parentSpanId,omitempty"`
	ServiceName  string `json:"serviceName,omitempty" db:"serviceName,omitempty"`
}

type ServiceDependency struct {
	Parent    string `json:"parent,omitempty"`
	Child     string `json:"child,omitempty"`
	CallCount int    `json:"callCount,omitempty"`
}

// operations

type OperationsArgs struct {
	ServiceName string
}

// endpoints

type TopEndpointsArgs struct {
	ServiceName string
	Start       *time.Time
	StartTime   string
	End         *time.Time
	EndTime     string
}

type TopEndpoints struct {
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

// rpc overview

type RpcOverview struct {
	Timestamp   int64   `json:"timestamp,omitempty" db:"timestamp,omitempty"`
	Time        string  `json:"time,omitempty" db:"time,omitempty"`
	RpcSystem   string  `json:"rpcSystem,omitempty" db:"rpcSystem,omitempty"`
	AvgDuration float32 `json:"avgDuration,omitempty" db:"avgDuration,omitempty"`
	NumCalls    int     `json:"numCalls,omitempty" db:"numCalls,omitempty"`
	CallRate    float32 `json:"callRate,omitempty" db:"callRate,omitempty"`
	NumErrors   int     `json:"numErrors" db:"numErrors"`
	ErrorRate   float32 `json:"errorRate" db:"errorRate"`
}

// db overview

// messaging overview

// spans

type SpansArgs struct {
	ServiceName   string
	OperationName string
	Kind          string
	Intervals     string
	Start         *time.Time
	End           *time.Time
	MinDuration   string
	MaxDuration   string
	Tags          []TagQuery
	Limit         int64
	Order         string
	Offset        int64
	BatchSize     int64
}

type Span struct {
	Columns []string        `json:"columns"`
	Events  [][]interface{} `json:"events"`
}

// span agg

type SpansAggregateArgs struct {
	ServiceName       string
	OperationName     string
	Kind              string
	Intervals         string
	Start             *time.Time
	End               *time.Time
	MinDuration       string
	MaxDuration       string
	Tags              []TagQuery
	GranOrigin        string
	GranPeriod        string
	StepSeconds       int
	Dimension         string
	AggregationOption string
}

type SpanAggregate struct {
	Timestamp int64   `json:"timestamp,omitempty" db:"timestamp"`
	Time      string  `json:"time,omitempty" db:"time"`
	Value     float32 `json:"value,omitempty" db:"value"`
}

// tags

type TagQuery struct {
	Key      string
	Value    string
	Operator string
}

type TagItem struct {
	TagKeys  string `json:"tagKeys" db:"tagKeys"`
	TagCount int    `json:"tagCount" db:"tagCount"`
}
