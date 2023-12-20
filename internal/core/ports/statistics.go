package ports

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type StatsPort interface {
	All(context.Context) ([]domain.Stats, error)
	Get(context.Context, string) (domain.Stats, error)
	Create(context.Context, domain.Stats) error
	Delete(context.Context, string) error
	GetStatsByLinkID(context.Context, string) ([]domain.Stats, error)
}
