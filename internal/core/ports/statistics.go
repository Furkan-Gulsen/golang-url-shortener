package ports

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type StatisticsDB interface {
	All(context.Context) ([]domain.Statistics, error)
	Get(context.Context, string) (*domain.Statistics, error)
	Create(context.Context, domain.Statistics) error
	Delete(context.Context, string) error
}
