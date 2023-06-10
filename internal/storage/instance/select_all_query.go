package instance

import (
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

func SelectByUserIDuery(tx storage.Transaction, userID string) ([]*Record, error) {
	records := make([]*Record, 0)

	err := tx.Select(
		&records,
		`SELECT * FROM instances WHERE user_id = $1 deleted_at IS NULL;`,
		userID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "instance: SelectAllQuery tx.Get error")
	}

	return records, nil
}
