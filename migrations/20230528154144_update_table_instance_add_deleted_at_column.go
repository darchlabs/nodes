package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upUpdateTableInstanceAddDeletedAtColumn, downUpdateTableInstanceAddDeletedAtColumn)
}

func upUpdateTableInstanceAddDeletedAtColumn(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		ALTER TABLE instances
		ADD COLUMN updated_at TIMESTAMPTZ,
		ADD COLUMN deleted_at TIMESTAMPTZ;
	`)
	if err != nil {
		return err
	}

	return nil
}

func downUpdateTableInstanceAddDeletedAtColumn(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
		ALTER TABLE instances
		DROP COLUMN updated_at,
		DROP COLUMN deleted_at;
	`)
	if err != nil {
		return err
	}

	return nil
}
