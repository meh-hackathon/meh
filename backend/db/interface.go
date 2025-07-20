package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Exec(sqlStr string, args ...any) (sql.Result, error) {
	return DB.Exec(sqlStr, args...)
}

func ExecContext(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	return DB.ExecContext(ctx, sqlStr, args...)
}

func Get(dest any, sqlStr string, args ...any) error {
	return DB.Get(dest, sqlStr, args...)
}

func GetContext(ctx context.Context, dest any, sqlStr string, args ...any) error {
	return DB.GetContext(ctx, dest, sqlStr, args...)
}

func Rebind(sqlStr string) string {
	return DB.Rebind(sqlStr)
}

func Select(dest any, sqlStr string, args ...any) error {
	return DB.Select(dest, sqlStr, args...)
}

func SelectContext(ctx context.Context, dest any, sqlStr string, args ...any) error {
	return DB.SelectContext(ctx, dest, sqlStr, args...)
}

func QueryRow(sqlStr string, args ...any) *sql.Row {
	return DB.QueryRow(sqlStr, args...)
}

func QueryRowContext(ctx context.Context, sqlStr string, args ...any) *sql.Row {
	return DB.QueryRowContext(ctx, sqlStr, args...)
}

func Query(sqlStr string, args ...any) (*sql.Rows, error) {
	return DB.Query(sqlStr, args...)
}

func QueryContext(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error) {
	return DB.QueryContext(ctx, sqlStr, args...)
}

func ExecTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := DB.BeginTxx(ctx, nil)
	if err != nil {
		return ErrTransaction.WithMessage("Failed to begin transaction").WithOrigin().WithCause(err)
	}
	defer tx.Rollback()
	return fn(tx)
}

func BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return DB.BeginTxx(ctx, opts)
}

func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}
