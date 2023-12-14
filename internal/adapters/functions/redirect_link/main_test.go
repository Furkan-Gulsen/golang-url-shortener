package main

import (
	"context"
	"testing"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/cache"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/tests/mock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func setupTest(shortLink string) (events.APIGatewayProxyResponse, error) {
	mockStore := mock.NewMockDynamoDBStore()
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	fillCache(cache, mockStore.Links)

	linkDomain := services.NewLinkDomain(mockStore, cache)
	apiHandler := handlers.NewAPIGatewayV2Handler(linkDomain)

	request := events.APIGatewayV2HTTPRequest{
		RawPath: "/" + shortLink,
	}
	response, err := apiHandler.Redirect(context.Background(), request)

	return response, err
}

func fillCache(cache *cache.RedisCache, links map[string]domain.Link) error {
	for _, link := range links {
		err := cache.Set(context.Background(), link.Id, link.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetOriginalLink_Success(t *testing.T) {
	response, err := setupTest("testid1")

	assert.NoError(t, err)
	assert.Equal(t, 301, response.StatusCode)

	location, ok := response.Headers["Location"]
	assert.True(t, ok)
	assert.Equal(t, "https://example.com/link1", location)
}

func TestGetOriginalLink_NotFound(t *testing.T) {
	response, err := setupTest("nonexistentid")
	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
	assert.Contains(t, response.Body, "Link not found")
}
