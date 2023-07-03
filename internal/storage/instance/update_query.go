package instance

import (
	"database/sql"
	"time"

	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

type UpdateQueryInput struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	Network     *string    `db:"network"`
	Environment *string    `db:"environment"`
	Name        *string    `db:"name"`
	ServiceURL  *string    `db:"service_url"`
	Artifacts   *Artifacts `json:"artifacts"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func UpdateQuery(tx storage.Transaction, input *UpdateQueryInput) error {
	_, err := tx.Exec(`
		UPDATE instances
		SET
			network = COALESCE($2, network),
			environment = COALESCE($3, environment),
			name = COALESCE($4, name),
			service_url = COALESCE($5, service_url),
			artifacts = COALESCE($6, artifacts),
			created_at = COALESCE($7, created_at),
			updated_at = COALESCE($8, updated_at),
			deleted_at = COALESCE($9, deleted_at)
		WHERE
			id = $1
		AND
			user_id = $10;`,
		input.ID,
		input.Network,
		input.Environment,
		input.Name,
		input.ServiceURL,
		input.Artifacts,
		input.CreatedAt,
		input.UpdatedAt,
		input.DeletedAt,
		input.UserID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return errors.Wrap(err, "instance: UpdateQuery tx.Exec error")
	}

	return nil
}
