package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func (h *ApiGatewayV2Handler) Redirect(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	pathSegments := strings.Split(req.RawPath, "/")
	if len(pathSegments) < 2 {
		return ClientError(http.StatusBadRequest, "Invalid URL path")
	}

	shortLinkKey := pathSegments[len(pathSegments)-1]
	longLink, err := h.link.GetOriginalURL(ctx, shortLinkKey)
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
