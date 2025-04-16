package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/observability/backend/internal/services/reader"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
)

type Metrics struct {
	reader *reader.Reader
	parser *RequestParser
}

func (m *Metrics) GetMetrics(w http.ResponseWriter, r *http.Request) {
	query, err := m.parser.ParseGetMetricsRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Metrics.GetMetrics", "failed to parse request: %v", err))
		return
	}

	result, err := m.reader.Metrics(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Metrics.GetMetrics", "failed to retrieve metrics: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewMetricsHandler(reader *reader.Reader, parser *RequestParser) *Metrics {
	m := &Metrics{
		reader: reader,
		parser: parser,
	}

	return m
}
