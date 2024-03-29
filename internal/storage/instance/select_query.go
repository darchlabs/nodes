package instance

import (
	"database/sql"

	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("instance: not found error")

type SelectQueryInput struct {
	ID string
}

func SelectQuery(tx storage.Transaction, input *SelectQueryInput) (*Record, error) {
	var record Record
	err := tx.Get(
		&record,
		`SELECT * FROM instances WHERE id = $1;`,
		input.ID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "instance: SelectAllQuery tx.Get error")
	}
	return &record, nil
}

// TODO: move to a different file
type SelectByUserIDQueryInput struct {
	ID     string
	UserID string
}

func SelectByUserIDQuery(tx storage.Transaction, input *SelectByUserIDQueryInput) (*Record, error) {
	var record Record
	err := tx.Get(
		&record,
		`SELECT * FROM instances WHERE id = $1 AND user_id = $2;`,
		input.ID,
		input.UserID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "instance: SelectAllQuery tx.Get error")
	}
	return &record, nil
}
