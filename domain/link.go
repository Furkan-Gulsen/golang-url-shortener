package domain

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/types"
)

type Link struct {
	db    types.DB
	cache types.Cache
}

func NewLinkDomain(d types.DB, c types.Cache) *Link {
	return &Link{db: d, cache: c}
}

func (service *Link) GetAllLinksFromDB(ctx context.Context) (*[]types.Link, error) {
	links, err := service.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all short URLs: %w", err)
	}
	return &links, nil
}

func (service *Link) GetOriginalURL(ctx context.Context, shortLinkKey string) (*string, error) {
	link, err := service.cache.Get(ctx, shortLinkKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get short URL for identifier '%s': %w", shortLinkKey, err)
	}
	return &link, nil
}

func (service *Link) Create(ctx context.Context, link types.Link) error {
	if err := service.cache.Set(ctx, link.Id, link.OriginalURL); err != nil {
		return fmt.Errorf("failed to set short URL for identifier '%s': %w", link.Id, err)
	}
	if err := service.db.Create(ctx, link); err != nil {
		return fmt.Errorf("failed to create short URL: %w", err)
	}
	return nil
}

func (service *Link) Delete(ctx context.Context, short string) error {
	if err := service.db.Delete(ctx, short); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", short, err)
	}
	if err := service.cache.Delete(ctx, short); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", short, err)
	}
	return nil
}
