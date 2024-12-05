package postgres

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
	"backend/internal/storage"
	"log"
	"fmt"
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