package domain

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/types"
)

type Link struct {
	db types.DB
}

func NewLinkDomain(d types.DB) *Link {
	return &Link{db: d}
}

func (service *Link) All(ctx context.Context) (*types.Link, error) {
	links, err := service.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all short URLs: %w", err)
	}
	return &links, nil
}

func (service *Link) Get(ctx context.Context, short string) (*types.Link, error) {
	link, err := service.db.Get(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("failed to get short URL for identifier '%s': %w", short, err)
	}
	return link, nil
}

func (service *Link) Create(ctx context.Context, link types.Link) error {
	if err := service.db.Create(ctx, link); err != nil {
		return fmt.Errorf("failed to create short URL: %w", err)
	}
	return nil
}

func (service *Link) Delete(ctx context.Context, short string) error {
	if err := service.db.Delete(ctx, short); err != nil {
		return fmt.Errorf("failed to delete short URL for identifier '%s': %w", short, err)
	}
	return nil
}
