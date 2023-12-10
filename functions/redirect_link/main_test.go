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

func setupTest(shortLink string) (events.APIGatewayProxyResponse, error) {
	mockStore := mock.NewMockDynamoDBStore()
	linkDomain := domain.NewLinkDomain(mockStore)
	apiHandler := handlers.NewAPIGatewayV2Handler(linkDomain)

	request := events.APIGatewayV2HTTPRequest{
		RawPath: "/" + shortLink,
	}
	response, err := apiHandler.GetOriginalLink(context.Background(), request)

	return response, err
}

func TestGetOriginalLink_Success(t *testing.T) {
	response, err := setupTest("testid1")

	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var link types.Link
	err = json.Unmarshal([]byte(response.Body), &link)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/link1", link.OriginalURL)
}

func TestGetOriginalLink_NotFound(t *testing.T) {
	response, err := setupTest("nonexistentid")

	assert.NoError(t, err)
	assert.Equal(t, 404, response.StatusCode)
	assert.Contains(t, response.Body, "Link not found")
}
