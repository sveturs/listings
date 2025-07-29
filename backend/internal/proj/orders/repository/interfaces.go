package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// TxBeginner интерфейс для начала транзакции
type TxBeginner interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// TxExecutor интерфейс для выполнения запросов в транзакции
type TxExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// DB интерфейс объединяющий возможности работы с БД
type DB interface {
	TxBeginner
	TxExecutor
}

// WithTx вспомогательная функция для выполнения операций в транзакции
func WithTx(ctx context.Context, db TxBeginner, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}