package ports

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type DB interface {
	All(context.Context) ([]domain.Link, error)
	Get(context.Context, string) (*domain.Link, error)
	Create(context.Context, domain.Link) error
	Delete(context.Context, string) error
}
