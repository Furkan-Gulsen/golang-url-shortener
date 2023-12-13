package mock

import (
	"context"
	"errors"
	"time"
)

type MockRedisCache struct {
	Store map[string]string
	TTL   map[string]time.Time
}

func NewMockRedisCache() *MockRedisCache {
	return &MockRedisCache{
		Store: make(map[string]string),
		TTL:   make(map[string]time.Time),
	}
}

func (m *MockRedisCache) Set(ctx context.Context, key string, val string) error {
	m.Store[key] = val
	m.TTL[key] = time.Now().Add(time.Minute)
	return nil
}

func (m *MockRedisCache) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.Store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	if time.Now().After(m.TTL[key]) {
		delete(m.Store, key)
		delete(m.TTL, key)
		return "", errors.New("key expired")
	}
	return val, nil
}

func (m *MockRedisCache) Delete(ctx context.Context, key string) error {
	_, ok := m.Store[key]
	if !ok {
		return errors.New("key not found")
	}
	delete(m.Store, key)
	delete(m.TTL, key)
	return nil
}
