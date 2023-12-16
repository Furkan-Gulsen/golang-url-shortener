package handlers

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type StatisticsFunctionHandler struct {
	linkService *services.LinkService
	statistics  *services.StatisticsService
}

func NewStatisticsFunctionHandler(l *services.LinkService, s *services.StatisticsService) *StatisticsFunctionHandler {
	return &StatisticsFunctionHandler{linkService: l, statistics: s}
}

func (h *StatisticsFunctionHandler) Statistics(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World!",
	}, nil

}
