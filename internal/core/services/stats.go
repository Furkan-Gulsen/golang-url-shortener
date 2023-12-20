package services

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/ports"
)

type StatsService struct {
	port  ports.StatsPort
	cache ports.Cache
}

func NewStatsService(p ports.StatsPort, c ports.Cache) *StatsService {
	return &StatsService{port: p, cache: c}
}

func (service *StatsService) All(ctx context.Context) ([]domain.Stats, error) {
	stats, err := service.port.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stats: %w", err)
	}
	return stats, nil
}

func (service *StatsService) Get(ctx context.Context, statsID string) (domain.Stats, error) {
	stats, err := service.port.Get(ctx, statsID)
	if err != nil {
		return domain.Stats{}, fmt.Errorf("failed to get stats for identifier '%s': %w", statsID, err)
	}
	return stats, nil
}

func (service *StatsService) Delete(ctx context.Context, linkID string) error {
	if err := service.port.Delete(ctx, linkID); err != nil {
		return fmt.Errorf("failed to delete stats for identifier '%s': %w", linkID, err)
	}
	return nil
}

func (service *StatsService) Create(ctx context.Context, data domain.Stats) error {
	if err := service.port.Create(ctx, data); err != nil {
		return fmt.Errorf("failed to create stats: %w", err)
	}
	return nil
}

func (service *StatsService) GetStatsByLinkID(ctx context.Context, linkID string) ([]domain.Stats, error) {
	stats, err := service.port.GetStatsByLinkID(ctx, linkID)
	if err != nil {
		return []domain.Stats{}, fmt.Errorf("failed to get stats for identifier '%s': %w", linkID, err)
	}
	return stats, nil
}
