package lib

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/go-sql-driver/mysql"
)

const (
	dbName = "wedding"

	openErr = "error opening mysql connection: %w"
)

type DBClient struct {
	ctx context.Context
	log *zap.Logger

	Db *sql.DB
}

func NewDBClient(ctx context.Context, cfg *Config, log *zap.Logger) (*DBClient, error) {
	log = log.Named("DBClient")

	username, err := cfg.Get(DBUsername)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	password, err := cfg.Get(DBPass)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	host, err := cfg.Get(DBHost)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	port, err := cfg.Get(DBPort)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}

	log = log.With(
		zap.String("Username", username),
		zap.String("Host", host),
		zap.String("Port", port),
		zap.String("Database", dbName),
	)

	config := &mysql.Config{
		User:      username,
		Passwd:    password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%v", host, port),
		DBName:    dbName,
		ParseTime: true,
	}

	log.Debug("Dialing mysql DB")
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	return &DBClient{ctx: ctx, log: log, Db: db}, nil
}

func (d *DBClient) CheckConnection() error {
	d.log.Debug("Pinging DB for health check...")
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
