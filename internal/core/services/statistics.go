package services

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/ports"
)

type StatisticsService struct {
	port  ports.StatisticsPort
	cache ports.Cache
}

func NewStatisticsService(p ports.StatisticsPort, c ports.Cache) *StatisticsService {
	return &StatisticsService{port: p, cache: c}
}

func (service *StatisticsService) Create(ctx context.Context, data domain.Statistics) error {
	if err := service.port.Create(ctx, data); err != nil {
		return fmt.Errorf("failed to create statistics: %w", err)
	}
	return nil
}

func (service *StatisticsService) All(ctx context.Context, linkID string) (*domain.Statistics, error) {
	statistics, err := service.port.Get(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for identifier '%s': %w", linkID, err)
	}
	return statistics, nil
}

func (service *StatisticsService) Delete(ctx context.Context, linkID string) error {
	if err := service.port.Delete(ctx, linkID); err != nil {
		return fmt.Errorf("failed to delete statistics for identifier '%s': %w", linkID, err)
	}
	return nil
}

func (service *StatisticsService) Get(ctx context.Context, linkID string) (*domain.Statistics, error) {
	statistics, err := service.port.Get(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for identifier '%s': %w", linkID, err)
	}
	return statistics, nil
}
