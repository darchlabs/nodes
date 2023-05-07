package instance

import (
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

type SelectQueryInput struct {
	ID string
}

func SelectQuery(tx storage.Transaction, input *SelectQueryInput) (*Record, error) {
	var record Record
	err := tx.Get(&record, `SELECT * FROM instances WHERE id = $1;`, input.ID)
	if err != nil {
		return nil, errors.Wrap(err, "instance: SelectQuery tx.Get error")
	}
	return &record, nil
}
