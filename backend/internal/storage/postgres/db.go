//backend/internal/storage/postgres/db.go
package postgres

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
	"backend/internal/storage"
	"log"
	"fmt"
    "github.com/jackc/pgx/v5" 
)

var _ storage.Storage = (*Database)(nil) // проверяем что Database реализует интерфейс Storage

type Database struct {
    pool *pgxpool.Pool
}

func NewDatabase(dbURL string) (*Database, error) {
    log.Printf("Connecting to database with URL: %s", dbURL)
    
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("error creating connection pool: %w", err)
    }

    // Проверяем подключение
    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("could not ping database: %w", err)
    }

    log.Printf("Successfully connected to database")
    
    return &Database{pool: pool}, nil
}

func (db *Database) Close() {
    if db.pool != nil {
        db.pool.Close()
    }
}

func (db *Database) Ping(ctx context.Context) error {
    return db.pool.Ping(ctx)
}
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

// Обновленный метод Query
func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (storage.Rows, error) {
    rows, err := db.pool.Query(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    return &RowsWrapper{rows: rows}, nil
}

// QueryRow остается без изменений
func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) storage.Row {
    return db.pool.QueryRow(ctx, sql, args...)
}