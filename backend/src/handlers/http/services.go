package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type Services struct {
	reader *reader.Reader
}

// get services

// get serviceMap deps

// get services list
func (s *Services) GetServicesList(w http.ResponseWriter, r *http.Request) {
	result, err := s.reader.ServicesList(context.TODO())
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Services.GetServicesList", "failed to retrieve services list: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewServicesHandler(reader *reader.Reader) *Services {
	s := &Services{
		reader: reader,
	}

	return s
}
