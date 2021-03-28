package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()

	DSN string
	Now func() time.Time
}

func NewDB() *DB {
	db := &DB{
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())

	return db
}

func (db *DB) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	// Connect to the database.
	if db.db, err = sql.Open("mysql", db.DSN); err != nil {
		return err
	}

	//ping to see if database connection is working
	if err := db.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx: tx,
		db: db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil

}