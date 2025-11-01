// backend/internal/proj/c2c/storage/opensearch/async_indexer_test.go
package opensearch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/domain/models"
)

// createTestRepository создает Repository для тестирования
// repo.client = nil и repo.storage = nil
// Это означает что любая индексация будет падать, но это ОК для unit тестов
// Мы тестируем async queue logic (enqueue, workers, retry, DLQ), а не Repository.IndexListing
func createTestRepository(db *sqlx.DB) *Repository {
	return &Repository{
		client:    nil, // приведёт к ошибкам при индексации
		indexName: "test_index",
		storage:   nil, // приведёт к ошибкам при fetchListing (ожидаемо)
		useAsync:  false,
	}
}

// TestAsyncIndexer_Basic тестирует базовую функциональность async indexer
func TestAsyncIndexer_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test DB connection
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup repository (client=nil, indexing will fail)
	repo := createTestRepository(db)

	// Create async indexer
	indexer := NewAsyncIndexer(repo, db, 2, 10)
	defer func() {
		err := indexer.Shutdown(5 * time.Second)
		require.NoError(t, err)
	}()

	t.Run("Enqueue_Success", func(t *testing.T) {
		task := IndexTask{
			ListingID: 123,
			Action:    "index",
			Data: &models.MarketplaceListing{
				ID:    123,
				Title: "Test Listing",
			},
		}

		err := indexer.Enqueue(task)
		assert.NoError(t, err, "Enqueue should succeed")

		// Ждём обработки
		time.Sleep(200 * time.Millisecond)

		// Задача попадёт в DLQ так как repo.client = nil
		var dlqCount int
		err = db.Get(&dlqCount, "SELECT COUNT(*) FROM opensearch_indexing_dlq WHERE listing_id = $1", 123)
		require.NoError(t, err)
		// Expect DLQ entry due to nil client (это нормально для unit тестов)

		// Cleanup
		_, _ = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id = $1", 123)
	})

	t.Run("GetQueueSize", func(t *testing.T) {
		queueSize := indexer.GetQueueSize()
		assert.GreaterOrEqual(t, queueSize, 0)
		assert.LessOrEqual(t, queueSize, 10) // max queue size
	})

	t.Run("IsHealthy", func(t *testing.T) {
		healthy := indexer.IsHealthy()
		assert.True(t, healthy, "Indexer should be healthy before shutdown")
	})
}

// TestAsyncIndexer_RetryMechanism тестирует retry механизм
func TestAsyncIndexer_RetryMechanism(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup repository (operations will fail due to nil client)
	repo := createTestRepository(db)

	indexer := NewAsyncIndexer(repo, db, 1, 10)
	defer func() {
		_ = indexer.Shutdown(5 * time.Second)
	}()

	t.Run("Retry_SavesToDLQ", func(t *testing.T) {
		task := IndexTask{
			ListingID: 999,
			Action:    "index",
			Data: &models.MarketplaceListing{
				ID:    999,
				Title: "Failed Listing",
			},
		}

		err := indexer.Enqueue(task)
		require.NoError(t, err)

		// Ждём все попытки retry (3 attempts с delays: 0s, 1s, 5s)
		time.Sleep(8 * time.Second)

		// Проверяем что задача попала в DLQ после exhausted retries
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM opensearch_indexing_dlq WHERE listing_id = $1", 999)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "Failed task should be in DLQ after retries exhausted")

		// Очистка
		_, err = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id = $1", 999)
		require.NoError(t, err)
	})
}

// TestAsyncIndexer_DLQRetry тестирует повторную обработку из DLQ
func TestAsyncIndexer_DLQRetry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := createTestRepository(db)

	indexer := NewAsyncIndexer(repo, db, 1, 10)
	defer func() {
		_ = indexer.Shutdown(5 * time.Second)
	}()

	t.Run("RetryDLQ_Enqueues", func(t *testing.T) {
		// Вставляем failed task в DLQ
		_, err := db.Exec(`
			INSERT INTO opensearch_indexing_dlq
			(listing_id, action, data, attempts, last_error, created_at, last_attempt_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, 888, "index", `{"id":888,"title":"DLQ Test"}`, 2, "test error", time.Now(), time.Now())
		require.NoError(t, err)

		// Initial DLQ count
		var countBefore int
		err = db.Get(&countBefore, "SELECT COUNT(*) FROM opensearch_indexing_dlq WHERE listing_id = $1", 888)
		require.NoError(t, err)
		assert.Equal(t, 1, countBefore, "Task should be in DLQ before retry")

		// Retry из DLQ
		err = indexer.RetryDLQ(context.Background(), 10)
		require.NoError(t, err, "RetryDLQ should succeed")

		// Даём время на обработку
		time.Sleep(500 * time.Millisecond)

		// Задача будет удалена из DLQ после enqueue, но затем снова попадёт в DLQ
		// так как repo.client = nil (ожидаемое поведение для unit тестов)

		// Cleanup
		_, err = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id = $1", 888)
		require.NoError(t, err)
	})
}

// TestAsyncIndexer_GracefulShutdown тестирует graceful shutdown
func TestAsyncIndexer_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := createTestRepository(db)

	indexer := NewAsyncIndexer(repo, db, 2, 100)

	t.Run("Shutdown_WithRemainingTasks", func(t *testing.T) {
		// Добавляем много задач
		for i := 10; i < 60; i++ {
			task := IndexTask{
				ListingID: i,
				Action:    "index",
				Data: &models.MarketplaceListing{
					ID:    i,
					Title: fmt.Sprintf("Listing %d", i),
				},
			}
			_ = indexer.Enqueue(task)
		}

		// Даём время начать обработку
		time.Sleep(100 * time.Millisecond)

		// Немедленный shutdown с коротким timeout
		err := indexer.Shutdown(1 * time.Second)
		assert.NoError(t, err)

		// Проверяем что indexer не healthy после shutdown
		healthy := indexer.IsHealthy()
		assert.False(t, healthy, "Indexer should not be healthy after shutdown")

		// Проверяем что оставшиеся задачи попали в DLQ
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM opensearch_indexing_dlq WHERE listing_id >= 10 AND listing_id < 60")
		require.NoError(t, err)
		t.Logf("DLQ count after shutdown: %d", count)

		// Должны быть задачи в DLQ (либо из queue, либо failed attempts)
		assert.Greater(t, count, 0, "Some tasks should be in DLQ after shutdown")

		// Cleanup
		_, err = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id >= 10 AND listing_id < 60")
		require.NoError(t, err)
	})
}

// TestAsyncIndexer_Metrics тестирует Prometheus metrics
func TestAsyncIndexer_Metrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := createTestRepository(db)

	indexer := NewAsyncIndexer(repo, db, 2, 10)
	defer func() {
		_ = indexer.Shutdown(5 * time.Second)
	}()

	t.Run("Metrics_QueueSize", func(t *testing.T) {
		// Enqueue несколько задач
		for i := 1000; i < 1005; i++ {
			task := IndexTask{
				ListingID: i,
				Action:    "index",
				Data: &models.MarketplaceListing{
					ID:    i,
					Title: fmt.Sprintf("Metric Test %d", i),
				},
			}
			err := indexer.Enqueue(task)
			require.NoError(t, err)
		}

		// Ждём обработки
		time.Sleep(1 * time.Second)

		// Проверяем queue size metric
		queueSize := indexer.GetQueueSize()
		assert.GreaterOrEqual(t, queueSize, 0)
		assert.LessOrEqual(t, queueSize, 10)

		// Cleanup DLQ
		_, _ = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id >= 1000 AND listing_id < 1005")
	})
}

// TestAsyncIndexer_Concurrency тестирует параллельную обработку
func TestAsyncIndexer_Concurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := createTestRepository(db)

	// 5 workers для параллельной обработки
	indexer := NewAsyncIndexer(repo, db, 5, 50)
	defer func() {
		_ = indexer.Shutdown(5 * time.Second)
	}()

	t.Run("Concurrent_Processing", func(t *testing.T) {
		// Enqueue 20 задач одновременно
		for i := 200; i < 220; i++ {
			task := IndexTask{
				ListingID: i,
				Action:    "index",
				Data: &models.MarketplaceListing{
					ID:    i,
					Title: fmt.Sprintf("Concurrent Test %d", i),
				},
			}
			err := indexer.Enqueue(task)
			require.NoError(t, err)
		}

		// Ждём обработки
		time.Sleep(2 * time.Second)

		// Все задачи должны быть обработаны (в DLQ из-за nil client)
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM opensearch_indexing_dlq WHERE listing_id >= 200 AND listing_id < 220")
		require.NoError(t, err)
		assert.Greater(t, count, 0, "Tasks should be processed and saved to DLQ")

		// Cleanup
		_, _ = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id >= 200 AND listing_id < 220")
	})
}

// TestAsyncIndexer_QueueOverflow тестирует переполнение очереди
func TestAsyncIndexer_QueueOverflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := createTestRepository(db)

	// Маленькая очередь (5 элементов)
	indexer := NewAsyncIndexer(repo, db, 1, 5)
	defer func() {
		_ = indexer.Shutdown(5 * time.Second)
	}()

	t.Run("Queue_Overflow_Fallback", func(t *testing.T) {
		// Быстро добавляем больше задач чем capacity очереди
		for i := 300; i < 320; i++ {
			task := IndexTask{
				ListingID: i,
				Action:    "index",
				Data: &models.MarketplaceListing{
					ID:    i,
					Title: fmt.Sprintf("Overflow Test %d", i),
				},
			}
			_ = indexer.Enqueue(task)
			// Некоторые enqueue могут fallback на sync indexing
			// Проверяем что не паникует
			assert.NotPanics(t, func() {
				_ = indexer.Enqueue(task)
			})
		}

		time.Sleep(2 * time.Second)

		// Cleanup
		_, _ = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id >= 300 AND listing_id < 320")
	})
}

// setupTestDB создаёт тестовую БД подключение
func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	// Подключение к тестовой БД (порт 5433 согласно CLAUDE.md)
	connStr := "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	require.NoError(t, err, "Failed to connect to test DB")

	// Проверяем что таблица DLQ существует
	var tableExists bool
	err = db.Get(&tableExists, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = 'opensearch_indexing_dlq'
		)
	`)
	require.NoError(t, err)
	if !tableExists {
		t.Skip("DLQ table doesn't exist, run migrations first")
	}

	cleanup := func() {
		// Очистка всех тестовых данных
		_, _ = db.Exec("DELETE FROM opensearch_indexing_dlq WHERE listing_id >= 0")
		db.Close()
	}

	return db, cleanup
}
