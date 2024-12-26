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

func (p *RequestParser) ParseGetServiceOverviewRequest(r *http.Request) (*reader.OverviewArgs, error) {
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

	fmt.Printf("STEP %+v", stepInt)

	fmt.Printf("SERVICE %v", serviceName)

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
