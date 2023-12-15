package ports

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type LinkDB interface {
	GetAllLinks(context.Context) ([]domain.Link, error)
	GetLink(context.Context, string) (*domain.Link, error)
	CreateLink(context.Context, domain.Link) error
	DeleteLink(context.Context, string) error

	CreateStatistics(context.Context, domain.Statistics) error
	DeleteStatistics(context.Context, string) error
}
