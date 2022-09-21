package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const driver = "mysql"

type StorageClient struct {
	db *sql.DB
}

func InitConnection(opts ConnectionOptions) *StorageClient {
	client := new(StorageClient)
	return client.Connect(opts)
}

func NewStorageClient(db *sql.DB) StorageClient {
	return StorageClient{
		db: db,
	}
}

// Connect establishes a connection with remote server.
func (c *StorageClient) Connect(opts ConnectionOptions) *StorageClient {
	chain := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", opts.User, opts.Password, opts.Host, opts.Port, opts.Schema)

	db, err := sql.Open(driver, chain)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(opts.ConnMaxOpen)
	db.SetMaxIdleConns(opts.ConnMaxIdle)
	db.SetConnMaxLifetime(opts.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	return &StorageClient{db: db}
}

// Close closes the connection with remote server.
func (c *StorageClient) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Exec prepares a SQL query and executes it. Usually used for modification.
func (c *StorageClient) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := c.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.ExecContext(ctx, args...)
}

// Query prepares a SQL query and executes it. Usually used for select.
func (c *StorageClient) Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := c.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.QueryContext(ctx, args...)
}

// QueryRow prepares a SQL query and executes it. Usually used for select at least one row.
func (c *StorageClient) QueryRow(ctx context.Context, sql string, args ...interface{}) *sql.Row {
	return c.db.QueryRowContext(ctx, sql, args)
}

// BeginTx starts a transaction for the given context
func (c *StorageClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}
