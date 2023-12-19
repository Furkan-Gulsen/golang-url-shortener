package handlers

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type DeleteFunctionHandler struct {
	statsService *services.StatsService
	linkService  *services.LinkService
}

func NewDeleteFunctionHandler(l *services.LinkService, s *services.StatsService) *DeleteFunctionHandler {
	return &DeleteFunctionHandler{linkService: l, statsService: s}
}

func (s *DeleteFunctionHandler) Delete(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]

	err := s.linkService.Delete(ctx, id)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	err = s.statsService.Delete(ctx, id)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{StatusCode: 204}, nil
}
