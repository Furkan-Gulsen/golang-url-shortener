package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
)

type RedirectFunctionHandler struct {
	linkService *services.LinkService
}

func NewRedirectFunctionHandler(l *services.LinkService) *RedirectFunctionHandler {
	return &RedirectFunctionHandler{linkService: l}
}

func (h *RedirectFunctionHandler) Redirect(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	pathSegments := strings.Split(req.RawPath, "/")
	if len(pathSegments) < 2 {
		return ClientError(http.StatusBadRequest, "Invalid URL path")
	}

	shortLinkKey := pathSegments[len(pathSegments)-1]
	longLink, err := h.linkService.GetOriginalURL(ctx, shortLinkKey)
	if err != nil || *longLink == "" {
		return ClientError(http.StatusNotFound, "Link not found")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Location": *longLink,
		},
	}, nil
}
