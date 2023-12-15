package mock

import (
	"context"
	"errors"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type MockDynamoDBStore struct {
	Links      map[string]domain.Link
	Statistics map[string]domain.Statistics
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

func (m *MockDynamoDBStore) GetAllLinks(ctx context.Context) ([]domain.Link, error) {
	var links []domain.Link
	for _, link := range m.Links {
		links = append(links, link)
	}
	return links, nil
}

func (m *MockDynamoDBStore) GetLink(ctx context.Context, id string) (*domain.Link, error) {
	if link, ok := m.Links[id]; ok {
		return &link, nil
	}
	return nil, errors.New("link not found")
}

func (m *MockDynamoDBStore) CreateLink(ctx context.Context, link domain.Link) error {
	if _, ok := m.Links[link.Id]; ok {
		return errors.New("link already exists")
	}
	m.Links[link.Id] = link
	return nil
}

func (m *MockDynamoDBStore) DeleteLink(ctx context.Context, id string) error {
	if _, ok := m.Links[id]; !ok {
		return errors.New("link not found")
	}
	delete(m.Links, id)
	return nil
}

func (m *MockDynamoDBStore) CreateStatistics(ctx context.Context, statistics domain.Statistics) error {
	if _, ok := m.Statistics[statistics.Id]; ok {
		return errors.New("statistics already exists")
	}
	m.Statistics[statistics.Id] = statistics
	return nil
}

func (m *MockDynamoDBStore) DeleteStatistics(ctx context.Context, id string) error {
	if _, ok := m.Statistics[id]; !ok {
		return errors.New("statistics not found")
	}
	delete(m.Statistics, id)
	return nil
}
