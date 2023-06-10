package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAlterTableInstancesAddUserId, downAlterTableInstancesAddUserId)
}

func upAlterTableInstancesAddUserId(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		ALTER TABLE instances ADD COLUMN user_id NOT NULL;
	`)
	if err != nil {
		return err
	}
	return nil
}

func downAlterTableInstancesAddUserId(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
