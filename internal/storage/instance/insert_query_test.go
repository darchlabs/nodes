package instance

import (
	"testing"
	"time"

	"github.com/darchlabs/nodes/internal/test"
	"github.com/jaekwon/testify/require"
	"github.com/jmoiron/sqlx"
)

func Test_InsertQuery_Integration(t *testing.T) {
	test.GetTxCall(t, func(db *sqlx.Tx, _ []interface{}) {
		record := &Record{
			ID:          "example-id",
			Network:     "darchlabs",
			Environment: "mainnet",
			Name:        "master-node",
			ServiceURL:  "http://server.url:5492",
			Artifacts: &Artifacts{
				Deployments: []string{"darch", "node"},
			},
			CreatedAt: time.Now(),
		}

		err := InsertQuery(db, record)
		require.NoError(t, err)

		var exp Record
		err = db.Get(&exp, `SELECT * FROM instances WHERE id = $1;`, record.ID)

		require.NoError(t, err)
		require.Equal(t, record.ID, exp.ID)
		require.Equal(t, record.Network, exp.Network)
		require.Equal(t, record.Environment, exp.Environment)
		require.Equal(t, record.Name, exp.Name)
		require.Equal(t, record.ServiceURL, exp.ServiceURL)
		require.Equal(t, record.ServiceURL, exp.ServiceURL)
		require.Equal(t, record.Artifacts, exp.Artifacts)
	})
}
