package instance

import (
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/pkg/errors"
)

func InsertQuery(tx storage.Transaction, record *Record) error {
	_, err := tx.Exec(`
		INSERT INTO instances (
			id,
			network,
			environment,
			name,
			artifacts,
			service_url,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		record.ID,
		record.Network,
		record.Environment,
		record.Name,
		record.Artifacts,
		record.ServiceURL,
		record.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "instance: InsertQuery tx.Exec error")
	}

	return nil
}
