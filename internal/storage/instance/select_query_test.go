package instance

import (
	"testing"

	"github.com/darchlabs/nodes/internal/test"
	"github.com/jaekwon/testify/require"
	"github.com/jmoiron/sqlx"
)

func Test_SelectQuery_Integration(t *testing.T) {
	test.GetTxCall(t, func(tx *sqlx.Tx, _ []interface{}) {
		expectedID := "example-id"

		_, err := tx.Exec(`
			INSERT INTO instances (
				id,
				network,
				environment,
				name,
				service_url,
				artifacts,
				created_at
			) VALUES (
				$1,
				'test',
				'dev',
				'test-record',
				'http://server.url:5492',
				'{"Deployments":["deployment-1","deployment-2"],"Pods":["pod-1","pod-2"],"Services":["service-1","service-2"]}'::jsonb,
				now()
			);`,
			expectedID,
		)
		require.NoError(t, err)

		record, err := SelectQuery(tx, &SelectQueryInput{ID: expectedID})

		require.NoError(t, err)
		require.Equal(t, expectedID, record.ID)
	})
}
