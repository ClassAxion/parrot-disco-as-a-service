package database

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

func NewPg(url string) *Database {
	connConfig, err := pgx.ParseConfig(url)
	if err != nil {
		panic(fmt.Errorf("failed to parse database url: %w", err))
	}

	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sqlx.Open("pgx", connStr)

	if err != nil {
		panic(err.Error())
	}

	return &Database{db}
}

func NewMysql(url string) *Database {
	db, err := sqlx.Open("mysql", url)

	if err != nil {
		panic(err.Error())
	}

	return &Database{db}
}

func (db *Database) UpdateQuery(ctx context.Context, query string, args ...any) error {
	_, err := db.ExecContext(ctx, query, args...)

	return err
}
