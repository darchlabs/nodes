package instance

import (
	"testing"

	"github.com/darchlabs/nodes/internal/test"
	"github.com/jaekwon/testify/require"
	"github.com/jmoiron/sqlx"
)

func Test_SelectAllQuery_Integration(t *testing.T) {
	test.GetTxCall(t, func(tx *sqlx.Tx, _ []interface{}) {
		expectedID1 := "example-id-1"
		expectedID2 := "example-id-2"

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
				'test-record-1',
				'http://server.url:5492',
				'{"Deployments":["deployment-1","deployment-2"],"Pods":["pod-1","pod-2"],"Services":["service-1","service-2"]}'::jsonb,
				now()
			);`,
			expectedID1,
		)
		require.NoError(t, err)

		_, err = tx.Exec(`
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
				'test-record-2',
				'http://server.url:5492',
				'{"Deployments":["deployment-1","deployment-2"],"Pods":["pod-1","pod-2"],"Services":["service-1","service-2"]}'::jsonb,
				now()
			);`,
			expectedID2,
		)
		require.NoError(t, err)

		records, err := SelectAllQuery(tx)

		require.NoError(t, err)
		require.Equal(t, 2, len(records))
	})
}
