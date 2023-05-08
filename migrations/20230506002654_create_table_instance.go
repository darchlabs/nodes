package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateTableInstances, downCreateTableInstances)
}

func upCreateTableInstances(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS instances (
			id TEXT PRIMARY KEY NOT NULL,
			network TEXT NOT NULL,
			environment TEXT NOT NULL,
			name TEXT UNIQUE NOT NULL,
			service_url TEXT NOT NULL,
			artifacts JSONB,
			created_at timestamptz NOT NULL
		);`)
	return err
}

func downCreateTableInstances(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
