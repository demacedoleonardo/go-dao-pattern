package mysql

import (
	"context"
	"database/sql"
	"time"
)

type (
	ConnectionOptions struct {
		Host            string
		User            string
		Password        string
		Schema          string
		RegisterName    string
		Port            int
		ConnMaxOpen     int
		ConnMaxIdle     int
		ConnMaxLifetime time.Duration
	}

	Client interface {
		Close() error
		Connect(opts ConnectionOptions) *StorageClient
		BeginTx(ctx context.Context, ops *sql.TxOptions) (*sql.Tx, error)
		Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
		Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	}
)
