package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/w-h-a/observability/backend/internal/clients/traces"
	"github.com/w-h-a/observability/backend/internal/handlers"
	"github.com/w-h-a/observability/backend/internal/services/reader"
)

type RequestParser struct{}

func (p *RequestParser) ParseGetServicesRequest(ctx context.Context, r *http.Request) (*reader.ServicesArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	serviceArgs := &reader.ServicesArgs{
		Start:     startTime,
		StartTime: startTime.Format(time.RFC3339Nano),
		End:       endTime,
		EndTime:   endTime.Format(time.RFC3339Nano),
		Period:    int(endTime.Unix() - startTime.Unix()),
	}

	return serviceArgs, nil
}

func (p *RequestParser) ParseGetTagsRequest(ctx context.Context, r *http.Request) (*reader.TagsArgs, error) {
	serviceName := r.URL.Query().Get("service")
	if len(serviceName) == 0 {
		return nil, errors.New("service param missing in query")
	}

	tagsArgs := &reader.TagsArgs{
		ServiceName: serviceName,
	}

	return tagsArgs, nil
}

func (p *RequestParser) ParseGetOperationsRequest(ctx context.Context, r *http.Request) (*reader.OperationsArgs, error) {
	serviceName := r.URL.Query().Get("service")
	if len(serviceName) == 0 {
		return nil, errors.New("service param missing in query")
	}

	operationsArgs := &reader.OperationsArgs{
		ServiceName: serviceName,
	}

	return operationsArgs, nil
}

func (p *RequestParser) ParseGetEndpointsRequest(ctx context.Context, r *http.Request) (*reader.EndpointsArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	serviceName := r.URL.Query().Get("service")
	if len(serviceName) == 0 {
		return nil, errors.New("service param missing in query")
	}

	topOperationsArgs := &reader.EndpointsArgs{
		ServiceName: serviceName,
		Start:       startTime,
		StartTime:   startTime.Format(time.RFC3339Nano),
		End:         endTime,
		EndTime:     endTime.Format(time.RFC3339Nano),
	}

	return topOperationsArgs, nil
}

func (p *RequestParser) ParseGetOverviewRequest(ctx context.Context, r *http.Request) (*reader.OverviewArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	stepStr := r.URL.Query().Get("step")
	if len(stepStr) == 0 {
		return nil, errors.New("step param missing in query")
	}

	stepInt, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, err
	}

	if stepInt < 60 {
		return nil, errors.New("step param is less than 60")
	}

	serviceName := r.URL.Query().Get("service")
	if len(serviceName) == 0 {
		return nil, errors.New("service param missing in query")
	}

	serviceOverviewArgs := &reader.OverviewArgs{
		ServiceName: serviceName,
		Start:       startTime,
		StartTime:   startTime.Format(time.RFC3339Nano),
		End:         endTime,
		EndTime:     endTime.Format(time.RFC3339Nano),
		StepSeconds: stepInt,
	}

	return serviceOverviewArgs, nil
}

func (p *RequestParser) ParseGetTracesRequest(ctx context.Context, r *http.Request) (*reader.TracesArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	tracesArgs := &reader.TracesArgs{
		Start:     startTime,
		StartTime: startTime.Format(time.RFC3339Nano),
		End:       endTime,
		EndTime:   endTime.Format(time.RFC3339Nano),
	}

	serviceName := r.URL.Query().Get("service")
	if len(serviceName) != 0 {
		tracesArgs.ServiceName = serviceName
	}

	traceId := r.URL.Query().Get("traceId")
	if len(traceId) != 0 {
		tracesArgs.TraceId = traceId
	}

	return tracesArgs, nil
}

func (p *RequestParser) ParseGetSpansRequest(ctx context.Context, r *http.Request) (*reader.SpansArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	spansArgs := &reader.SpansArgs{
		Start:     startTime,
		StartTime: startTime.Format(time.RFC3339Nano),
		End:       endTime,
		EndTime:   endTime.Format(time.RFC3339Nano),
	}

	serviceName := r.URL.Query().Get("service")
	if len(serviceName) != 0 {
		spansArgs.ServiceName = serviceName
	}

	spanName := r.URL.Query().Get("name")
	if len(spanName) != 0 {
		spansArgs.Name = spanName
	}

	spanKind := r.URL.Query().Get("kind")
	if len(spanKind) != 0 {
		spansArgs.Kind = spanKind
	}

	minDuration, err := p.parseTimestamp(r.URL.Query().Get("minDuration"))
	if err == nil {
		spansArgs.MinDuration = *minDuration
	}

	maxDuration, err := p.parseTimestamp(r.URL.Query().Get("maxDuration"))
	if err == nil {
		spansArgs.MaxDuration = *maxDuration
	}

	tagQueries, err := p.parseTagQueries(r.URL.Query().Get("tags"))
	if err != nil {
		return nil, err
	}
	if len(tagQueries) != 0 {
		spansArgs.TagQueries = tagQueries
	}

	// TODO: order, limit, offset

	return spansArgs, nil
}

func (p *RequestParser) ParseGetAggregatedSpansRequest(ctx context.Context, r *http.Request) (*reader.AggregatedSpansArgs, error) {
	startTime, err := p.parseTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	stepStr := r.URL.Query().Get("step")
	if len(stepStr) == 0 {
		return nil, errors.New("step param missing in query")
	}

	stepInt, err := strconv.Atoi(stepStr)
	if err != nil {
		return nil, err
	}

	if stepInt < 60 {
		return nil, errors.New("step param is less than 60")
	}

	dimension := r.URL.Query().Get("dimension")
	if len(dimension) == 0 {
		return nil, errors.New("dimension param missing in query")
	}
	if !slices.Contains(handlers.SupportedRequestedDimensions, dimension) {
		return nil, fmt.Errorf("dimension %s is not supported in query", dimension)
	}

	aggregationOption := r.URL.Query().Get("aggregation")
	if len(aggregationOption) == 0 {
		return nil, errors.New("aggregation param missing in query")
	}
	if !slices.Contains(handlers.SupportedRequestedAggregations[dimension], aggregationOption) {
		return nil, fmt.Errorf("aggregation %s is not supported in query with dimension %s", aggregationOption, dimension)
	}

	spanAggregatesArgs := &reader.AggregatedSpansArgs{
		Dimension:         dimension,
		AggregationOption: aggregationOption,
		StepSeconds:       stepInt,
		Start:             startTime,
		StartTime:         startTime.Format(time.RFC3339Nano),
		End:               endTime,
		EndTime:           endTime.Format(time.RFC3339Nano),
	}

	serviceName := r.URL.Query().Get("service")
	if len(serviceName) != 0 {
		spanAggregatesArgs.ServiceName = serviceName
	}

	spanName := r.URL.Query().Get("name")
	if len(spanName) != 0 {
		spanAggregatesArgs.Name = spanName
	}

	spanKind := r.URL.Query().Get("kind")
	if len(spanKind) != 0 {
		spanAggregatesArgs.Kind = spanKind
	}

	minDuration, err := p.parseTimestamp(r.URL.Query().Get("minDuration"))
	if err == nil {
		spanAggregatesArgs.MinDuration = *minDuration
	}

	maxDuration, err := p.parseTimestamp(r.URL.Query().Get("maxDuration"))
	if err == nil {
		spanAggregatesArgs.MaxDuration = *maxDuration
	}

	tagQueries, err := p.parseTagQueries(r.URL.Query().Get("tags"))
	if err != nil {
		return nil, err
	}
	if len(tagQueries) != 0 {
		spanAggregatesArgs.TagQueries = tagQueries
	}

	return spanAggregatesArgs, nil
}

func (p *RequestParser) ParseGetMetricsRequest(ctx context.Context, r *http.Request) (*reader.MetricsArgs, error) {
	start, err := p.parseMetricsTime(r.URL.Query().Get("start"))
	if err != nil {
		return nil, err
	}

	end, err := p.parseMetricsTime(r.URL.Query().Get("end"))
	if err != nil {
		return nil, err
	}

	step, err := p.parseMetricsDuration(r.URL.Query().Get("step"))
	if err != nil {
		return nil, err
	}

	if (*end).Sub(*start)/step > 10000 {
		return nil, errors.New("exceeded max resolution of 10,000 data points")
	}

	dimension := r.URL.Query().Get("dimension")
	if len(dimension) == 0 {
		return nil, errors.New("dimension param missing in query")
	}
	if !slices.Contains(handlers.SupportedMetrics, dimension) {
		return nil, fmt.Errorf("dimension %s is not supported in query", dimension)
	}

	serviceName := r.URL.Query().Get("service")
	if slices.Contains(handlers.RequiredMetricAttributes[dimension], "service") && len(serviceName) == 0 {
		return nil, fmt.Errorf("service param is required with dimension %s", dimension)
	}

	metricsArgs := &reader.MetricsArgs{
		Dimension:   dimension,
		Step:        step,
		Start:       start,
		End:         end,
		ServiceName: serviceName,
	}

	return metricsArgs, nil
}

func (p *RequestParser) parseTime(timeStr string) (*time.Time, error) {
	if len(timeStr) == 0 {
		return nil, fmt.Errorf("time param missing in query")
	}

	timeUnix, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("time param is not in correct timestamp format")
	}

	timeFmt := time.Unix(timeUnix, 0)

	return &timeFmt, nil
}

func (p *RequestParser) parseMetricsTime(timeStr string) (*time.Time, error) {
	if len(timeStr) == 0 {
		return nil, fmt.Errorf("time param missing in query")
	}

	if timeUnix, err := strconv.ParseFloat(timeStr, 64); err == nil {
		s, ns := math.Modf(timeUnix)
		timeFmt := time.Unix(int64(s), int64(ns*float64(time.Second)))
		return &timeFmt, nil
	}

	if timeFmt, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return &timeFmt, nil
	}

	return &time.Time{}, nil
}

func (p *RequestParser) parseMetricsDuration(durStr string) (time.Duration, error) {
	if len(durStr) == 0 {
		return 0, errors.New("step param missing in query")
	}

	if d, err := strconv.ParseFloat(durStr, 64); err == nil {
		ts := d * float64(time.Second)
		if ts > float64(math.MaxInt64) || ts < float64(60) {
			return 0, fmt.Errorf("%q overflows the legal limit", durStr)
		}
		return time.Duration(ts), nil
	}

	return 0, fmt.Errorf("failed to parse %q as a duration", durStr)
}

func (p *RequestParser) parseTimestamp(timeStr string) (*string, error) {
	if len(timeStr) == 0 {
		return nil, fmt.Errorf("timestamp param missing in query")
	}

	return &timeStr, nil
}

func (p *RequestParser) parseTagQueries(tagsStr string) ([]traces.TagQuery, error) {
	tagQueries := []traces.TagQuery{}

	if len(tagsStr) == 0 {
		return tagQueries, nil
	}

	if err := json.Unmarshal([]byte(tagsStr), &tagQueries); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %v", err)
	}

	for _, tagQuery := range tagQueries {
		if !slices.Contains(handlers.SupportedRequestedTagOperators, tagQuery.Operator) {
			return nil, fmt.Errorf("tag operator %s is not supported in tag query", tagQuery.Operator)
		}
	}

	return tagQueries, nil
}
