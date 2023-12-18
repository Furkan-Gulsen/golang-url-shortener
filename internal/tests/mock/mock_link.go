package mock

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type MockLinkRepo struct {
	Links []domain.Link
	Stats []domain.Stats
}

func NewMockLinkRepo() *MockLinkRepo {
	return &MockLinkRepo{
		Links: MockLinkData,
		Stats: MockStatsData,
	}
}

func (m *MockLinkRepo) All(ctx context.Context) ([]domain.Link, error) {
	return m.Links, nil
}

func (m *MockLinkRepo) Get(ctx context.Context, id string) (domain.Link, error) {
	for _, link := range m.Links {
		if link.Id == id {
			return link, nil
		}
	}

	return domain.Link{}, nil
}

func (m *MockLinkRepo) Create(ctx context.Context, link domain.Link) error {
	m.Links = append(m.Links, link)
	return nil
}

func (m *MockLinkRepo) Delete(ctx context.Context, id string) error {
	for i, link := range m.Links {
		if link.Id == id {
			m.Links = append(m.Links[:i], m.Links[i+1:]...)
			return nil
		}
	}

	return nil
}
