package mock

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
)

type MockStatisticsRepo struct {
	Statistics map[string]domain.Statistics
}

func NewMockStatisticsRepo() *MockStatisticsRepo {
	return &MockStatisticsRepo{
		Statistics: map[string]domain.Statistics{
			"testid1": {Id: "abcdefg1", ClickCount: 0, Platform: domain.PlatformUnknown, LinkID: "testid1"},
			"testid2": {Id: "abcdefg2", ClickCount: 0, Platform: domain.PlatformUnknown, LinkID: "testid2"},
			"testid3": {Id: "abcdefg3", ClickCount: 0, Platform: domain.PlatformUnknown, LinkID: "testid3"},
		},
	}
}

func (m *MockStatisticsRepo) Get(ctx context.Context, id string) (*domain.Statistics, error) {
	if statistics, ok := m.Statistics[id]; ok {
		return &statistics, nil
	}
	return nil, nil
}

func (m *MockStatisticsRepo) Create(ctx context.Context, statistics domain.Statistics) error {
	if _, ok := m.Statistics[statistics.Id]; ok {
		return nil
	}
	m.Statistics[statistics.Id] = statistics
	return nil
}

func (m *MockStatisticsRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.Statistics[id]; !ok {
		return nil
	}
	delete(m.Statistics, id)
	return nil
}
