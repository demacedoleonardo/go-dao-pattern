package db

import (
	"context"
	"database/sql"
	"go-dao-pattern/pkg/metrics"
	"go-dao-pattern/pkg/storage/mysql"
)

type Action string

const (
	SELECT Action = "SELECT"
	INSERT Action = "INSERT"
	UPDATE Action = "UPDATE"
	DELETE Action = "DELETE"
)

func (o Action) String() string {
	return string(o)
}

// ExecQuery executes a query and return rows
func ExecQuery(ctx context.Context, client mysql.Client, resource, query string, args ...interface{}) (*sql.Rows, error) {
	var (
		rows *sql.Rows
		err  error
	)

	metrics.StartStoreSegment(func() error {
		rows, err = client.Query(ctx, query, args...)
		return err
	},
		metrics.WithAction(SELECT.String()),
		metrics.WithResource(resource),
		metrics.WithContext(ctx))

	return rows, err
}

// ExecQueryRow executes a query and return single row
func ExecQueryRow(ctx context.Context, client mysql.Client, resource, query string, args ...interface{}) *sql.Row {
	var (
		row *sql.Row
	)

	metrics.StartStoreSegment(func() error {
		row = client.QueryRow(ctx, query, args...)
		return nil
	},
		metrics.WithAction(SELECT.String()),
		metrics.WithResource(resource),
		metrics.WithContext(ctx))

	return row
}

// ExecStatement executes a statement and return result
func ExecStatement(ctx context.Context, c mysql.Client, a Action, resource, query string, args ...interface{}) (sql.Result, error) {
	var (
		result sql.Result
		err    error
	)

	metrics.StartStoreSegment(func() error {
		result, err = c.Exec(ctx, query, args...)
		return err
	},
		metrics.WithAction(a.String()),
		metrics.WithResource(resource),
		metrics.WithContext(ctx))

	return result, err
}

// ExecStatementWithTx executes a transactional statement
func ExecStatementWithTx(ctx context.Context, a Action, sqlTx *sql.Tx, resource, query string, args ...interface{}) (sql.Result, error) {
	var err error
	var result sql.Result

	metrics.StartStoreSegment(func() error {
		result, err = sqlTx.ExecContext(ctx, query, args...)
		return err
	},
		metrics.WithAction(a.String()),
		metrics.WithResource(resource),
		metrics.WithContext(ctx))

	return result, err
}
