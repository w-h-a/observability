package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type Traces struct {
	reader *reader.Reader
	parser *RequestParser
}

func (t *Traces) GetTraces(w http.ResponseWriter, r *http.Request) {
	query, err := t.parser.ParseGetTracesRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Traces.GetTraces", "failed to parse request: %v", err))
		return
	}

	result, err := t.reader.Traces(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Traces.GetTraces", "failed to retrieve traces: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewTracesHandler(reader *reader.Reader, parser *RequestParser) *Traces {
	t := &Traces{
		reader: reader,
		parser: parser,
	}

	return t
}
