package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type StatsFunctionHandler struct {
	statsService *services.StatsService
	linkService  *services.LinkService
}

func NewStatsFunctionHandler(l *services.LinkService, s *services.StatsService) *StatsFunctionHandler {
	return &StatsFunctionHandler{linkService: l, statsService: s}
}

func (s *StatsFunctionHandler) Stats(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	links, err := s.linkService.GetAll(ctx)
	if err != nil {
		return ServerError(err)
	}

	for i, link := range links {
		stats, err := s.statsService.GetStatsByLinkID(ctx, link.Id)
		if err != nil {
			log.Printf("Error getting stats for link '%s': %v", link.Id, err)
			continue
		}
		links[i].Stats = stats
	}

	jsonResponse, err := json.Marshal(links)
	if err != nil {
		return ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResponse),
	}, nil
}
