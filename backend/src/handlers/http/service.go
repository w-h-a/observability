package http

import (
	"context"
	"net/http"

	"github.com/w-h-a/pkg/utils/errorutils"
	"github.com/w-h-a/pkg/utils/httputils"
	"github.com/w-h-a/trace-blame/backend/src/handlers/http/utils"
	"github.com/w-h-a/trace-blame/backend/src/services/reader"
)

type Service struct {
	reader *reader.Reader
	parser *utils.RequestParser
}

func (s *Service) GetOperations(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetOperationsRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Service.GetOperations", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.Operations(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Service.GetOperations", "failed to retrieve service operations: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Service) GetEndpoints(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetEndpointsRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Service.GetEndpoints", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.Endpoints(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Service.GetEndpoints", "failed to retrieve top endpoints: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Service) GetServiceOverview(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetOverviewRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Service.GetServiceOverview", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.ServiceOverview(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Service.GetServiceOverview", "failed to retrieve service overview: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func (s *Service) GetTags(w http.ResponseWriter, r *http.Request) {
	query, err := s.parser.ParseGetTagsRequest(context.TODO(), r)
	if err != nil {
		httputils.ErrResponse(w, errorutils.BadRequest("Service.GetTags", "failed to parse request: %v", err))
		return
	}

	result, err := s.reader.Tags(context.TODO(), query)
	if err != nil {
		httputils.ErrResponse(w, errorutils.InternalServerError("Service.GetTags", "failed to retrieve tags: %v", err))
		return
	}

	httputils.OkResponse(w, result)
}

func NewServiceHandler(reader *reader.Reader, parser *utils.RequestParser) *Service {
	s := &Service{
		reader: reader,
		parser: parser,
	}

	return s
}
