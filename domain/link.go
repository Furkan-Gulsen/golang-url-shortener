package domain

import (
	"context"
	"fmt"

	"github.com/Furkan-Gulsen/golang-url-shortener/types"
)

type Link struct {
	db types.DB
}

func NewLink(db types.DB) *Link {
	return &Link{db}
}

func (s *Link) All(ctx context.Context) ([]types.Link, error) {
	link, err := s.db.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all short urls: %w", err)
	}
	return link, nil
}

func (s *Link) Get(ctx context.Context, short string) (types.Link, error) {
	link, err := s.db.Get(ctx, short)
	if err != nil {
		return types.Link{}, fmt.Errorf("error getting short url: %w", err)
	}
	return link, nil
}

func (s *Link) Create(ctx context.Context, short types.Link) error {
	err := s.db.Create(ctx, short)
	if err != nil {
		return fmt.Errorf("error creating short url: %w", err)
	}
	return nil
}

func (s *Link) Delete(ctx context.Context, short string) error {
	err := s.db.Delete(ctx, short)
	if err != nil {
		return fmt.Errorf("error deleting short url: %w", err)
	}
	return nil
}
