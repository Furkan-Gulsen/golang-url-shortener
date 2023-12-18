package handlers

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

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
