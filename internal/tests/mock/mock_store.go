package mock

import (
	"context"
	"errors"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type MockDynamoDBStore struct {
	Links map[string]domain.Link
}

func NewMockDynamoDBStore() *MockDynamoDBStore {
	return &MockDynamoDBStore{
		Links: map[string]domain.Link{
			"testid1": {Id: "testid1", OriginalURL: "https://example.com/link1"},
			"testid2": {Id: "testid2", OriginalURL: "https://example.com/link2"},
			"testid3": {Id: "testid3", OriginalURL: "https://example.com/link3"},
		},
	}
}

func (m *MockDynamoDBStore) All(ctx context.Context) ([]domain.Link, error) {
	var links []domain.Link
	for _, link := range m.Links {
		links = append(links, link)
	}
	return links, nil
}

func (m *MockDynamoDBStore) Get(ctx context.Context, id string) (*domain.Link, error) {
	if link, ok := m.Links[id]; ok {
		return &link, nil
	}
	return nil, errors.New("link not found")
}

func (m *MockDynamoDBStore) Create(ctx context.Context, link domain.Link) error {
	if _, ok := m.Links[link.Id]; ok {
		return errors.New("link already exists")
	}
	m.Links[link.Id] = link
	return nil
}

func (m *MockDynamoDBStore) Delete(ctx context.Context, id string) error {
	if _, ok := m.Links[id]; !ok {
		return errors.New("link not found")
	}
	delete(m.Links, id)
	return nil
}
