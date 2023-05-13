package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/darchlabs/nodes/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

const (
	defaultTestDBDriver = "postgres"
	defaultTestDBDSN    = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(fmt.Sprintf("%s/../../../", b))
)

func getDB() (*sqlx.DB, error) {
	driver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("POSTGRES_DSN")
	migrations := os.Getenv("POSTGRES_MIGRATIONS_DIR")

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	db.Exec("drop database if exists postgres")
	db.Exec("create postgres")

	db.SetConnMaxLifetime(-1)
	err = goose.Up(db.DB, fmt.Sprintf("%s/%s", basepath, migrations))
	if err != nil {
		return nil, err
	}

	return db, nil
}

type TestDB struct {
	tx *sqlx.Tx
}

func (tdb *TestDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return tdb.tx, nil
}

func (tdb *TestDB) Exec(query string, params ...interface{}) (sql.Result, error) {
	return tdb.tx.Exec(query, params)
}

func (tdb *TestDB) Query(query string, params ...interface{}) (*sql.Rows, error) {
	return tdb.tx.Query(query, params)
}

func (tdb *TestDB) Get(dest interface{}, q string, args ...interface{}) error {
	return tdb.tx.Get(dest, q, args)
}

func (tdb *TestDB) Select(dest interface{}, q string, args ...interface{}) error {
	return tdb.tx.Select(dest, q, args)
}

func (tdb *TestDB) Rollback() error {
	return tdb.tx.Rollback()
}

func GetTxCall(t *testing.T, call func(tx *sqlx.Tx, testData []interface{})) {
	t.Helper()
	conn, err := getDB()
	require.NoError(t, err)

	ctx := context.Background()
	tx, err := conn.BeginTxx(ctx, nil)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	ctx = context.Background()
	tx1, err := conn.BeginTxx(ctx, nil)
	require.NoError(t, err)

	call(tx1, nil)

	tx1.Commit()

	ctx = context.Background()
	tx, err = conn.BeginTxx(ctx, nil)
	require.NoError(t, err)

	err = CleanDB(tx)
	require.NoError(t, err)

	tx.Commit()

	err = conn.Close()
	require.NoError(t, err)
}

func GetDBCall(t *testing.T, call func(db *sqlx.DB, testData []interface{})) {
	t.Helper()
	conn, err := getDB()
	require.NoError(t, err)

	CleanDBConn(t, conn)

	ctx := context.Background()
	tx, err := conn.BeginTxx(ctx, nil)
	require.NoError(t, err)

	tx.Commit()

	call(conn, nil)

	CleanDBConn(t, conn)

	conn.Close()
}

func CleanDBConn(t *testing.T, st *sqlx.DB) error {
	err := PrepareDeleteFromDB(st, "DELETE FROM instances CASCADE;")
	if err != nil {
		return err
	}

	return nil
}

func PrepareDeleteFrom(st *sqlx.Tx, query string) error {
	stmt, err := st.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func PrepareDeleteFromDB(st *sqlx.DB, query string) error {
	stmt, err := st.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func CleanDB(st *sqlx.Tx) (err error) {
	if st == nil {
		return nil
	}

	_, err = st.Exec("DELETE FROM instances;")
	if err != nil {
		return err
	}

	return nil
}
