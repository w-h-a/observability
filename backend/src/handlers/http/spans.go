package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/observability/backend/src/services/reader"
	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
)

type Spans struct {
	reader *reader.Reader
	parser *RequestParser
}

func (s *Spans) GetSpans(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetSpansRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Spans.GetSpans", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.Spans(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Spans.GetSpans", "failed to retrieve spans: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Spans) GetAggregatedSpans(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetAggregatedSpansRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Spans.GetAggregatedSpans", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.AggregatedSpans(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Spans.GetAggregatedSpans", "failed to retrieve aggregated spans: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewSpansHandler(reader *reader.Reader, parser *RequestParser) *Spans {
	s := &Spans{
		reader: reader,
		parser: parser,
	}

	return s
}
