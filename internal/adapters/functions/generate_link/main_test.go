package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/cache"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/tests/mock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func setupTest(body string) (events.APIGatewayProxyResponse, error) {
	mockStore := &mock.MockDynamoDBStore{
		Links: map[string]domain.Link{},
	}
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	linkDomain := services.NewLinkDomain(mockStore, cache)
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

	var link domain.Link
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
