package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/cache"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/tests/mock"
	"github.com/google/uuid"
)

func GetService() *services.LinkService {
	cache := cache.NewRedisCache("localhost:6379", "", 0)
	mockLinkRepo := mock.NewMockLinkRepo()

	linkService := services.NewLinkService(mockLinkRepo, cache)

	return linkService
}

func BenchmarkLinkServiceGetAll(b *testing.B) {
	service := GetService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetAll(ctx)
		if err != nil {
			b.Fatalf("Benchmark GetAll failed: %v", err)
		}
	}
}

func BenchmarkLinkServiceGetOriginalURL(b *testing.B) {
	service := GetService()
	ctx := context.Background()
	shortLinkKey := "testid2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetOriginalURL(ctx, shortLinkKey)
		if err != nil {
			b.Fatalf("Benchmark GetOriginalURL failed: %v", err)
		}
	}
}

// Benchmark Create
func BenchmarkLinkServiceCreate(b *testing.B) {
	service := GetService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newID := uuid.New().String()
		link := domain.Link{
			Id:          newID,
			OriginalURL: "https://example.com/" + newID,
			CreatedAt:   time.Now(),
		}

		err := service.Create(ctx, link)
		if err != nil {
			b.Fatalf("Benchmark Create failed: %v", err)
		}
	}
}
