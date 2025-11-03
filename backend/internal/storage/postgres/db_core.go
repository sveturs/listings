// backend/internal/storage/postgres/db_core.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	marketplaceStorage "backend/internal/proj/marketplace/storage"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Close закрывает все соединения с базой данных
func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
	if db.db != nil {
		if err := db.db.Close(); err != nil {
			// Логируем ошибку, но не прерываем выполнение закрытия
			_ = err // Explicitly ignore error
		}
	}
}

// GetSQLXDB возвращает sqlx.DB для использования в модулях, которые требуют sqlx
func (db *Database) GetSQLXDB() *sqlx.DB {
	// Если sqlxDB уже инициализирован, возвращаем его
	if db.sqlxDB != nil {
		return db.sqlxDB
	}
	// Создаем новый sqlx.DB из пула и СОХРАНЯЕМ его
	stdDB := stdlib.OpenDBFromPool(db.pool)
	db.sqlxDB = sqlx.NewDb(stdDB, "pgx")
	return db.sqlxDB
}

// FileStorage возвращает интерфейс файлового хранилища
func (db *Database) FileStorage() filestorage.FileStorageInterface {
	return db.fsStorage
}

// Marketplace возвращает интерфейс для работы с marketplace storage
func (db *Database) Marketplace() marketplaceStorage.MarketplaceStorage {
	return db.marketplaceStorage
}

// GetSQLDB returns the raw sql.DB connection
func (db *Database) GetSQLDB() *sql.DB {
	return db.db
}

// Ping проверяет соединение с базой данных
func (db *Database) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// pgxResult реализует sql.Result для pgx
type pgxResult struct {
	ct pgconn.CommandTag
}

func (r pgxResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId is not supported by PostgreSQL")
}

func (r pgxResult) RowsAffected() (int64, error) {
	return r.ct.RowsAffected(), nil
}

// Exec выполняет SQL запрос без возврата результата
func (db *Database) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	ct, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxResult{ct: ct}, nil
}

// RowsWrapper обертка для pgx.Rows для реализации storage.Rows
type RowsWrapper struct {
	rows pgx.Rows
}

func (r *RowsWrapper) Next() bool {
	return r.rows.Next()
}

func (r *RowsWrapper) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *RowsWrapper) Close() error {
	r.rows.Close()
	return nil
}

func (r *RowsWrapper) Err() error {
	return r.rows.Err()
}

// Query выполняет SQL запрос и возвращает строки результата
func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (storage.Rows, error) {
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &RowsWrapper{rows: rows}, nil
}

// QueryRow выполняет SQL запрос и возвращает одну строку результата
func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) storage.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

// QueryContext выполняет SQL запрос с контекстом и возвращает строки результата
func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

// QueryRowContext выполняет SQL запрос с контекстом и возвращает одну строку результата
func (db *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

// ExecContext выполняет SQL запрос с контекстом без возврата результата
func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// pgxTransaction реализует интерфейс Transaction для pgx
type pgxTransaction struct {
	tx pgx.Tx
}

// BeginTx начинает новую транзакцию
func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (storage.Transaction, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTransaction{tx: tx}, nil
}

// Реализация методов интерфейса Transaction
func (t *pgxTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ct, err := t.tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgxResult{ct: ct}, nil
}

func (t *pgxTransaction) Query(ctx context.Context, query string, args ...interface{}) (storage.Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &RowsWrapper{rows: rows}, nil
}

func (t *pgxTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) storage.Row {
	return t.tx.QueryRow(ctx, query, args...)
}

func (t *pgxTransaction) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t *pgxTransaction) Rollback() error {
	return t.tx.Rollback(context.Background())
}

// GetAttributeGroups возвращает репозиторий групп атрибутов
func (db *Database) GetAttributeGroups() AttributeGroupStorage {
	return db.attributeGroups
}

// Cart возвращает репозиторий корзин
func (db *Database) Cart() interface{} {
	if db.cartRepo != nil {
		return db.cartRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewCartRepository(db.pool)
}

// Order возвращает репозиторий заказов
func (db *Database) Order() interface{} {
	if db.orderRepo != nil {
		return db.orderRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewOrderRepository(db.pool)
}

// Inventory возвращает репозиторий инвентаря
func (db *Database) Inventory() interface{} {
	if db.inventoryRepo != nil {
		return db.inventoryRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewInventoryRepository(db.pool)
}

// MarketplaceOrder возвращает репозиторий заказов маркетплейса
func (db *Database) MarketplaceOrder() interface{} {
	if db.marketplaceOrderRepo != nil {
		return db.marketplaceOrderRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewMarketplaceOrderRepository(db.pool)
}

// StorefrontProductSearch возвращает репозиторий для поиска товаров витрин
// TODO: Migrate to marketplace microservice
func (db *Database) StorefrontProductSearch() interface{} {
	if db.productSearchRepo != nil {
		return db.productSearchRepo
	}
	// OpenSearch disabled after removing b2c
	return nil
}
