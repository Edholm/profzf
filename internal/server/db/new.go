package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	// sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

type CleanupFunc func() error

type Mode string

const (
	ModeReadOnly        Mode = "ro"
	ModeReadWriteCreate Mode = "rwc"
)

func NewQuerier(ctx context.Context, file string, mode Mode) (Querier, CleanupFunc, error) {
	ds := fmt.Sprintf("file:%s?cache=shared&mode=%s", file, mode)
	db, err := sql.Open("sqlite3", ds)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("sql.Open: %w", err)
	}
	// Execute DDL
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, db.Close, fmt.Errorf("db.ExecContext: %w", err)
	}
	return New(db), db.Close, err
}
