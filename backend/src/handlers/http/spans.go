package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type Spans struct {
	reader *reader.Reader
	parser *utils.RequestParser
}

func (s *Spans) GetSpansByTraceId(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetSpansByTraceRequest(r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Spans.GetSpansByTraceId", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.SpansByTrace(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Spans.GetSpansByTraceId", "failed to retrieve spans by traceId: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewSpansHandler(reader *reader.Reader, parser *utils.RequestParser) *Spans {
	s := &Spans{
		reader: reader,
		parser: parser,
	}

	return s
}
