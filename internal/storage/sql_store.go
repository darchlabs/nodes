package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Store is the database wrapper
type SQLStore struct {
	*sqlx.DB
}

func NewSQLStore(driver, dsn string) (*SQLStore, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// Maximum Idle Connections
	db.SetMaxIdleConns(20)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)

	return &SQLStore{db}, nil
}

func (st *SQLStore) BeginTx(ctx context.Context) (Transactioner, error) {
	tx, err := st.DB.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, errors.Wrap(err, "database: SQLStore.BeginTx st.BeginTxx error")
	}
	return tx, nil
}
