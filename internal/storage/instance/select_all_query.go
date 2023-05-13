package instance

import (
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

func SelectAllQuery(tx storage.Transaction) ([]*Record, error) {
	records := make([]*Record, 0)

	err := tx.Select(&records, `SELECT * FROM instances;`)
	if err != nil {
		return nil, errors.Wrap(err, "instance: SelectAllQuery tx.Get error")
	}

	return records, nil
}
