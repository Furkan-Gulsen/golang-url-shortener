package services

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/ports"
)

type StatisticsService struct {
	db    ports.StatisticsDB
	cache ports.Cache
}

func NewStatisticsService(d ports.StatisticsDB, c ports.Cache) *StatisticsService {
	return &StatisticsService{db: d, cache: c}
}

func (service *StatisticsService) Create(ctx context.Context, data domain.Statistics) error {
	if err := service.db.Create(ctx, data); err != nil {
		return fmt.Errorf("failed to create statistics: %w", err)
	}
	return nil
}

func (service *StatisticsService) All(ctx context.Context, linkID string) (*domain.Statistics, error) {
	statistics, err := service.db.Get(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for identifier '%s': %w", linkID, err)
	}
	return statistics, nil
}

func (service *StatisticsService) Delete(ctx context.Context, linkID string) error {
	if err := service.db.Delete(ctx, linkID); err != nil {
		return fmt.Errorf("failed to delete statistics for identifier '%s': %w", linkID, err)
	}
	return nil
}

// GET
func (service *StatisticsService) Get(ctx context.Context, linkID string) (*domain.Statistics, error) {
	statistics, err := service.db.Get(ctx, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics for identifier '%s': %w", linkID, err)
	}
	return statistics, nil
}
