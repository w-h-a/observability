package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type RequestParser struct{}

func (p *RequestParser) ParseGetServicesRequest(r *http.Request) (*reader.ServicesArgs, error) {
	startTime, err := p.parseTime("start", r)
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime("end", r)
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

func (p *RequestParser) ParseGetOperationsRequest(r *http.Request) (*reader.OperationsArgs, error) {
	serviceName := r.URL.Query().Get("service")
	if len(serviceName) == 0 {
		return nil, errors.New("service param missing in query")
	}

	operationsArgs := &reader.OperationsArgs{
		ServiceName: serviceName,
	}

	return operationsArgs, nil
}

func (p *RequestParser) ParseGetEndpointsRequest(r *http.Request) (*reader.EndpointsArgs, error) {
	startTime, err := p.parseTime("start", r)
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime("end", r)
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

func (p *RequestParser) ParseGetOverviewRequest(r *http.Request) (*reader.OverviewArgs, error) {
	startTime, err := p.parseTime("start", r)
	if err != nil {
		return nil, err
	}

	endTime, err := p.parseTime("end", r)
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

func (p *RequestParser) ParseGetSpansByTraceRequest(r *http.Request) (*reader.SpansByTraceIdArgs, error) {
	traceId := r.URL.Query().Get("traceId")
	if len(traceId) == 0 {
		return nil, errors.New("traceId param missing in query")
	}

	spansByTraceIdArgs := &reader.SpansByTraceIdArgs{
		TraceId: traceId,
	}

	return spansByTraceIdArgs, nil
}

func (p *RequestParser) parseTime(param string, r *http.Request) (*time.Time, error) {
	timeStr := r.URL.Query().Get(param)

	if len(timeStr) == 0 {
		return nil, fmt.Errorf("%s param missing in query", param)
	}

	timeUnix, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s param is not in correct timestamp format", param)
	}

	timeFmt := time.Unix(timeUnix, 0)

	return &timeFmt, nil
}
