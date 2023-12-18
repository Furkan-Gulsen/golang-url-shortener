package unit

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

func TestStatsTest(t *testing.T) {
	mockStatsRepo := mock.NewMockStatsRepo()
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	statsService := services.NewStatsService(mockStatsRepo, cache)

	mockLinkRepo := mock.NewMockLinkRepo()
	linkService := services.NewLinkService(mockLinkRepo, cache)

	apiHander := handlers.NewStatsFunctionHandler(linkService, statsService)

	t.Run("Stats Unit Test", func(t *testing.T) {
		request := events.APIGatewayV2HTTPRequest{
			RawPath: "/stats",
		}

		response, err := apiHander.Stats(context.Background(), request)
		if err != nil {
			t.Fatal(err)
		}

		var links []domain.Link
		err = json.Unmarshal([]byte(response.Body), &links)

		assert.Nil(t, err)
		assert.Equal(t, response.StatusCode, 200)
		assert.Equal(t, len(links), 3)

		assert.Equal(t, links[0].Stats.ClickCount, 0)
		assert.Equal(t, links[0].Stats.Platform, domain.PlatformUnknown)
		assert.Equal(t, links[0].Stats.LinkID, "testid1")

		assert.Equal(t, links[1].Stats.ClickCount, 1)
		assert.Equal(t, links[1].Stats.Platform, domain.PlatformInstagram)
		assert.Equal(t, links[1].Stats.LinkID, "testid2")

		assert.Equal(t, links[2].Stats.ClickCount, 2)
		assert.Equal(t, links[2].Stats.LinkID, "testid3")
		assert.Equal(t, links[2].Stats.Platform, domain.PlatformTwitter)

	})

}
