package mock

import (
	"context"
	"errors"

	"github.com/Furkan-Gulsen/golang-url-shortener/types"
)

type MockDynamoDBStore struct {
	Links map[string]types.Link
}

func NewMockDynamoDBStore() *MockDynamoDBStore {
	return &MockDynamoDBStore{
		Links: map[string]types.Link{
			"testid1": {Id: "testid1", OriginalURL: "https://example.com/link1"},
			"testid2": {Id: "testid2", OriginalURL: "https://example.com/link2"},
			"testid3": {Id: "testid3", OriginalURL: "https://example.com/link3"},
		},
	}
}

func (m *MockDynamoDBStore) All(ctx context.Context) ([]types.Link, error) {
	var links []types.Link
	for _, link := range m.Links {
		links = append(links, link)
	}
	return links, nil
}

func (m *MockDynamoDBStore) Get(ctx context.Context, id string) (*types.Link, error) {
	if link, ok := m.Links[id]; ok {
		return &link, nil
	}
	return nil, errors.New("link not found")
}

func (m *MockDynamoDBStore) Create(ctx context.Context, link types.Link) error {
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
