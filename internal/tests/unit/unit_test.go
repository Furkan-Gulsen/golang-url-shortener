package unit

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/cache"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/tests/mock"
)

func SetupTest() *handlers.ApiGatewayV2Handler {
	mockStore := mock.NewMockDynamoDBStore()
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	fillCache(cache, mockStore.Links)

	linkDomain := services.NewLinkDomain(mockStore, cache)
	apiHandler := handlers.NewAPIGatewayV2Handler(linkDomain)

	return apiHandler

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
