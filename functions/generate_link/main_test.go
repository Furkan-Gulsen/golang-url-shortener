package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/mock"
	"github.com/Furkan-Gulsen/golang-url-shortener/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func setupTest(body string) (events.APIGatewayProxyResponse, error) {
	mockStore := &mock.MockDynamoDBStore{
		Links: map[string]types.Link{},
	}
	linkDomain := domain.NewLinkDomain(mockStore)
	apiHandler := handlers.NewAPIGatewayV2Handler(linkDomain)

	request := events.APIGatewayV2HTTPRequest{Body: body}
	response, err := apiHandler.CreateShortLink(context.Background(), request)

	return response, err
}

func TestCreateShortLink_Success(t *testing.T) {
	body := `{"long": "https://example.com/abcdefg"}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var link types.Link
	err = json.Unmarshal([]byte(response.Body), &link)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/abcdefg", link.OriginalURL)
}

func TestCreateShortLink_EmptyString(t *testing.T) {
	body := `{"long": ""}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, "URL cannot be empty", response.Body)
}

func TestCreateShortLink_InvalidURL(t *testing.T) {
	body := `{"long": "invalid"}`
	response, err := setupTest(body)

	assert.NoError(t, err)
	assert.Equal(t, 400, response.StatusCode)
	assert.Equal(t, "URL must be at least 15 characters long", response.Body)
}
