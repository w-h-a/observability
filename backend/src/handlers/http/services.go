package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type Services struct {
	reader *reader.Reader
	parser *utils.RequestParser
}

func (s *Services) GetServices(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetServicesRequest(r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Services.GetServices", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.Services(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Services.GetServices", "failed to retrieve services: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Services) GetServicesList(w http.ResponseWriter, r *http.Request) {
	result, err := s.reader.ServicesList(context.TODO())
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Services.GetServicesList", "failed to retrieve services list: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Services) GetServiceDependencies(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetServicesRequest(r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Services.GetServiceDependencies", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.ServiceDependencies(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Services.GetServiceDependencies", "failed to retrieve service map: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewServicesHandler(reader *reader.Reader, parser *utils.RequestParser) *Services {
	s := &Services{
		reader: reader,
		parser: parser,
	}

	return s
}
