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
		expectedUserID := "example-user-id"

		_, err := tx.Exec(`
			INSERT INTO instances (
				id,
				user_id,
				network,
				environment,
				name,
				service_url,
				artifacts,
				created_at
			) VALUES (
				$1,
				$2
				'test',
				'dev',
				'test-record',
				'http://server.url:5492',
				'{"Deployments":["deployment-1","deployment-2"],"Pods":["pod-1","pod-2"],"Services":["service-1","service-2"]}'::jsonb,
				now()
			);`,
			expectedID,
			expectedUserID,
		)
		require.NoError(t, err)

		record, err := SelectQuery(tx, &SelectQueryInput{ID: expectedID, UserID: expectedUserID})

		require.NoError(t, err)
		require.Equal(t, expectedID, record.ID)
	})
}
