package instance

import (
	"testing"
	"time"

	"github.com/darchlabs/nodes/internal/test"
	"github.com/jaekwon/testify/require"
	"github.com/jmoiron/sqlx"
)

func Test_UpdateQuery_integration(t *testing.T) {
	test.GetTxCall(t, func(db *sqlx.Tx, _ []interface{}) {
		_, err := db.Exec(`
			INSERT INTO instances (
				id,
				network,
				environment,
				name,
				service_url,
				artifacts,
				created_at
			) VALUES (
				'test-id',
				'some-network',
				'mainnet',
				'test-node',
				'http://node.com/rpc',
				'{"foo": 123, "bar": "baz"}'::jsonb,
				now()
			);
		`)
		require.NoError(t, err)

		now := time.Now()
		err = UpdateQuery(db, &UpdateQueryInput{
			ID:        "test-id",
			DeletedAt: &now,
		})

		require.NoError(t, err)

		var record Record
		err = db.Get(&record, `SELECT * FROM instances WHERE id = 'test-id';`)

		require.NoError(t, err)
		require.NotNil(t, record.DeletedAt)
	})
}
