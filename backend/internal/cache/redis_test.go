package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*RedisCache, *miniredis.Miniredis) {
	// Create a mock Redis server
	mr, err := miniredis.Run()
	require.NoError(t, err)

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create Redis cache pointing to mock server
	cache, err := NewRedisCache(context.Background(), mr.Addr(), "", 0, 10, logger)
	require.NoError(t, err)

	return cache, mr
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		key     string
		value   interface{}
		ttl     time.Duration
		wantErr bool
	}{
		{
			name:    "Set and get simple struct",
			key:     "test_key",
			value:   testData{ID: 1, Name: "Test"},
			ttl:     time.Hour,
			wantErr: false,
		},
		{
			name:    "Set with zero TTL",
			key:     "test_key_no_ttl",
			value:   testData{ID: 2, Name: "NoTTL"},
			ttl:     0,
			wantErr: false,
		},
		{
			name:    "Set complex data",
			key:     "complex_key",
			value:   map[string]interface{}{"nested": map[string]string{"key": "value"}},
			ttl:     time.Hour,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := cache.Set(ctx, tt.key, tt.value, tt.ttl)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Get value
			var result interface{}
			if _, ok := tt.value.(testData); ok {
				result = &testData{}
			} else {
				result = &map[string]interface{}{}
			}

			err = cache.Get(ctx, tt.key, result)
			assert.NoError(t, err)

			// Compare values
			if td, ok := tt.value.(testData); ok {
				resultData := result.(*testData)
				assert.Equal(t, td.ID, resultData.ID)
				assert.Equal(t, td.Name, resultData.Name)
			}
		})
	}
}

func TestRedisCache_Delete(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Set a value
	key := "delete_test"
	value := map[string]string{"test": "value"}
	err := cache.Set(ctx, key, value, time.Hour)
	require.NoError(t, err)

	// Verify it exists
	var result map[string]string
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Delete the key
	err = cache.Delete(ctx, key)
	assert.NoError(t, err)

	// Verify it's gone
	err = cache.Get(ctx, key, &result)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestRedisCache_DeleteMultiple(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Set multiple values
	keys := []string{"key1", "key2", "key3"}
	for i, key := range keys {
		err := cache.Set(ctx, key, fmt.Sprintf("value%d", i), time.Hour)
		require.NoError(t, err)
	}

	// Delete all keys
	err := cache.Delete(ctx, keys...)
	assert.NoError(t, err)

	// Verify all are gone
	for _, key := range keys {
		var result string
		err := cache.Get(ctx, key, &result)
		assert.ErrorIs(t, err, ErrCacheMiss)
	}
}

func TestRedisCache_DeletePattern(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Set values with pattern
	testData := map[string]string{
		"categories:1": "cat1",
		"categories:2": "cat2",
		"categories:3": "cat3",
		"attributes:1": "attr1",
	}

	for key, value := range testData {
		err := cache.Set(ctx, key, value, time.Hour)
		require.NoError(t, err)
	}

	// Delete by pattern
	err := cache.DeletePattern(ctx, "categories:*")
	assert.NoError(t, err)

	// Verify categories are deleted
	for key := range testData {
		var result string
		err := cache.Get(ctx, key, &result)
		if key[:10] == "categories" {
			assert.ErrorIs(t, err, ErrCacheMiss)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestRedisCache_GetNonExistent(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Try to get non-existent key
	var result string
	err := cache.Get(ctx, "non_existent", &result)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestRedisCache_Exists(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Check non-existent key
	exists, err := cache.Exists(ctx, "non_existent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Set a key
	key := "exist_test"
	err = cache.Set(ctx, key, "value", time.Hour)
	require.NoError(t, err)

	// Check it exists
	exists, err = cache.Exists(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRedisCache_TTLExpiry(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Set value with short TTL
	key := "ttl_test"
	value := "expires_soon"
	ttl := 100 * time.Millisecond

	err := cache.Set(ctx, key, value, ttl)
	require.NoError(t, err)

	// Value should exist immediately
	var result string
	err = cache.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Fast forward time in miniredis
	mr.FastForward(200 * time.Millisecond)

	// Value should be expired
	err = cache.Get(ctx, key, &result)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestRedisCache_GetOrSet(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	type category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	key := "getorset_test"
	expectedValue := category{ID: 1, Name: "Test Category"}

	// First call - should call loader
	loaderCalled := false
	var result category
	err := cache.GetOrSet(ctx, key, &result, time.Hour, func() (interface{}, error) {
		loaderCalled = true
		return expectedValue, nil
	})

	assert.NoError(t, err)
	assert.True(t, loaderCalled)
	assert.Equal(t, expectedValue, result)

	// Wait a bit for async caching to complete
	time.Sleep(100 * time.Millisecond)

	// Second call - should get from cache
	loaderCalled = false
	var result2 category
	err = cache.GetOrSet(ctx, key, &result2, time.Hour, func() (interface{}, error) {
		loaderCalled = true
		return expectedValue, nil
	})

	assert.NoError(t, err)
	assert.False(t, loaderCalled)
	assert.Equal(t, expectedValue, result2)
}

func TestRedisCache_GetOrSetLoaderError(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	key := "getorset_error"
	expectedError := fmt.Errorf("loader error")

	var result string
	err := cache.GetOrSet(ctx, key, &result, time.Hour, func() (interface{}, error) {
		return nil, expectedError
	})

	assert.ErrorIs(t, err, expectedError)
}

func TestRedisCache_ConnectionFailure(t *testing.T) {
	logger := logrus.New()

	// Try to create cache with invalid connection
	_, err := NewRedisCache(context.Background(), "localhost:99999", "", 0, 10, logger)
	assert.Error(t, err)
}

func TestRedisCache_ConcurrentAccess(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Test concurrent set operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			key := fmt.Sprintf("concurrent_%d", n)
			value := fmt.Sprintf("value_%d", n)
			err := cache.Set(ctx, key, value, time.Hour)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values are set
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent_%d", i)
		expectedValue := fmt.Sprintf("value_%d", i)
		var result string
		err := cache.Get(ctx, key, &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	}
}

// TestRedisCache_GracefulDegradation tests graceful degradation when Redis is unavailable
func TestRedisCache_GracefulDegradation(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()

	ctx := context.Background()

	// Close the connection to simulate Redis being unavailable
	if err := cache.Close(); err != nil {
		t.Logf("Failed to close cache: %v", err)
	}

	t.Run("Set with closed connection", func(t *testing.T) {
		err := cache.Set(ctx, "test", "value", time.Hour)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "closed")
	})

	t.Run("Get with closed connection", func(t *testing.T) {
		var result string
		err := cache.Get(ctx, "test", &result)
		assert.Error(t, err)
		// Should not be ErrCacheMiss as this is a connection error
		assert.NotErrorIs(t, err, ErrCacheMiss)
	})

	t.Run("Delete with closed connection", func(t *testing.T) {
		err := cache.Delete(ctx, "test")
		assert.Error(t, err)
	})

	t.Run("Exists with closed connection", func(t *testing.T) {
		exists, err := cache.Exists(ctx, "test")
		assert.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("GetOrSet with closed connection - graceful fallback", func(t *testing.T) {
		type testData struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		var result testData
		expected := testData{ID: 1, Name: "Fallback Value"}

		// GetOrSet should return the loader result even when cache fails
		err := cache.GetOrSet(ctx, "test", &result, time.Hour, func() (interface{}, error) {
			return expected, nil
		})

		// Should succeed with data from loader
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

// TestRedisCache_MarshalUnmarshalErrors tests error handling for serialization issues
func TestRedisCache_MarshalUnmarshalErrors(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	t.Run("Marshal error - channel type", func(t *testing.T) {
		// Channels cannot be marshaled to JSON
		type unmarshalable struct {
			Ch chan int `json:"ch"`
		}

		value := unmarshalable{Ch: make(chan int)}
		err := cache.Set(ctx, "unmarshalable", value, time.Hour)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "json")
	})

	t.Run("Unmarshal error - type mismatch", func(t *testing.T) {
		// Set a string value
		err := cache.Set(ctx, "type_mismatch", "string value", time.Hour)
		require.NoError(t, err)

		// Try to get it as a struct
		type testStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		var result testStruct
		err = cache.Get(ctx, "type_mismatch", &result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshal")
	})

	t.Run("Invalid JSON in cache", func(t *testing.T) {
		// Manually set invalid JSON in Redis
		err := cache.client.Set(ctx, "invalid_json", "{invalid json", 0).Err()
		require.NoError(t, err)

		var result map[string]string
		err = cache.Get(ctx, "invalid_json", &result)
		assert.Error(t, err)
		// Check that it's a JSON parsing error
		assert.Contains(t, err.Error(), "invalid character")
	})
}

// TestRedisCache_DeleteEmptyKeys tests deleting with empty key list
func TestRedisCache_DeleteEmptyKeys(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	// Should not error when deleting zero keys
	err := cache.Delete(ctx)
	assert.NoError(t, err)
}

// TestRedisCache_Close tests proper connection closing
func TestRedisCache_Close(t *testing.T) {
	cache, mr := setupTestRedis(t)
	defer func() {
		mr.Close()
	}()

	ctx := context.Background()

	// Ensure connection works before closing
	err := cache.Set(ctx, "before_close", "value", time.Hour)
	require.NoError(t, err)

	// Close the connection
	err = cache.Close()
	assert.NoError(t, err)

	// Operations should fail after close
	err = cache.Set(ctx, "after_close", "value", time.Hour)
	assert.Error(t, err)
}

// TestRedisCache_GetOrSetCacheError tests GetOrSet behavior when cache has errors
func TestRedisCache_GetOrSetCacheError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create cache with valid miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer func() {
		mr.Close()
	}()

	cache, err := NewRedisCache(context.Background(), mr.Addr(), "", 0, 10, logger)
	require.NoError(t, err)

	ctx := context.Background()
	key := "test_cache_error"

	// Set up a key with invalid data that will cause unmarshal error
	err = cache.client.Set(ctx, key, "invalid json for struct", time.Hour).Err()
	require.NoError(t, err)

	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// GetOrSet should handle cache error and use loader
	var result testData
	expected := testData{ID: 1, Name: "Loaded Value"}
	loaderCalled := false

	err = cache.GetOrSet(ctx, key, &result, time.Hour, func() (interface{}, error) {
		loaderCalled = true
		return expected, nil
	})

	assert.NoError(t, err)
	assert.True(t, loaderCalled)
	assert.Equal(t, expected, result)
}

// Benchmarks

func BenchmarkRedisCache_Set(b *testing.B) {
	mr, err := miniredis.Run()
	require.NoError(b, err)
	defer func() {
		mr.Close()
	}()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cache, err := NewRedisCache(context.Background(), mr.Addr(), "", 0, 10, logger)
	require.NoError(b, err)
	defer func() {
		if err := cache.Close(); err != nil {
			b.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	data := testData{ID: 1, Name: "Benchmark User"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("benchmark:set:%d", i)
			_ = cache.Set(ctx, key, data, time.Hour)
			i++
		}
	})
}

func BenchmarkRedisCache_Get(b *testing.B) {
	mr, err := miniredis.Run()
	require.NoError(b, err)
	defer func() {
		mr.Close()
	}()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cache, err := NewRedisCache(context.Background(), mr.Addr(), "", 0, 10, logger)
	require.NoError(b, err)
	defer func() {
		if err := cache.Close(); err != nil {
			b.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	data := testData{ID: 1, Name: "Benchmark User"}

	// Pre-populate cache
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("benchmark:get:%d", i)
		_ = cache.Set(ctx, key, data, time.Hour)
	}

	var result testData
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("benchmark:get:%d", i%100)
			_ = cache.Get(ctx, key, &result)
			i++
		}
	})
}

func BenchmarkRedisCache_GetOrSet(b *testing.B) {
	mr, err := miniredis.Run()
	require.NoError(b, err)
	defer func() {
		mr.Close()
	}()

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cache, err := NewRedisCache(context.Background(), mr.Addr(), "", 0, 10, logger)
	require.NoError(b, err)
	defer func() {
		if err := cache.Close(); err != nil {
			b.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	data := testData{ID: 1, Name: "Benchmark User"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var result testData
			key := fmt.Sprintf("benchmark:getorset:%d", i%100)
			_ = cache.GetOrSet(ctx, key, &result, time.Hour, func() (interface{}, error) {
				return data, nil
			})
			i++
		}
	})
}
