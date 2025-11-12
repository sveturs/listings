package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockCacheRepository is a mock implementation of listings.CacheRepository interface
type MockCacheRepository struct {
	mock.Mock
}

// Get mocks getting a value from cache
func (m *MockCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Set mocks setting a value in cache
func (m *MockCacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

// Delete mocks deleting a value from cache
func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
