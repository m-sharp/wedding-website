package lib

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	dbName = "wedding"
)

type DBClient struct {
	ctx *context.Context
	Db  *sql.DB
}

func NewDBClient(ctx *context.Context, host, username, password string, port int) (*DBClient, error) {
	config := &mysql.Config{
		User:      username,
		Passwd:    password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%v", host, port),
		DBName:    dbName,
		ParseTime: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening mysql connection: %w", err)
	}
	return &DBClient{ctx: ctx, Db: db}, nil
}

//func (d *DBClient) Execute(ctx context.Context, query string, args ...interface{}) error {
//	if _, err := d.Db.ExecContext(ctx, query, args...); err != nil {
//		return errors.New(fmt.Sprintf("failed to execute query %q: %s", query, err))
//	}
//	return nil
//}
//
//func (d *DBClient) Query(ctx context.Context, query string, args ...interface{}) (int, error) {
//	rows, err := d.Db.QueryContext(ctx, query, args...)
//	if err != nil {
//		return 0, err
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			println("Error closing Rows: %s", err)
//		}
//	}(rows)
//	// ToDo - query row context vs query context? Either way, needs to be pluggable for use downstream
//	return 0, nil
//}

func (d *DBClient) CheckConnection() error {
	return d.Db.Ping()
}

type DBError struct {
	inner error
	query string
}

func NewDBError(query string, innerErr error) *DBError {
	return &DBError{inner: innerErr, query: query}
}

func (d *DBError) Error() string {
	return fmt.Sprintf("failed to execute query %q: %s", d.query, d.inner)
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		println("Error closing Rows: %s", err)
	}
}
