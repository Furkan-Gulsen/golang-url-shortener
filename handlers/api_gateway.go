package handlers

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/aws/aws-lambda-go/events"
)

type ApiGatewayV2Handler struct {
	link *domain.Link
}

func NewAPIGatewayV2Handler(l *domain.Link) *ApiGatewayV2Handler {
	return &ApiGatewayV2Handler{link: l}
}

func ClientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
	}, nil
}

func ServerError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

func IsValidLink(u string) bool {
	re := regexp.MustCompile(`^(http|https)://`)
	if !re.MatchString(u) {
		return false
	}

	parsedURL, err := url.ParseRequestURI(u)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}
