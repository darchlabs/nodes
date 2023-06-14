package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	instancedb "github.com/darchlabs/nodes/internal/storage/instance"
)

type nodeManager interface {
	DeployNewNode(*manager.CreateDeploymentOptions) (*manager.NodeInstance, error)
	DeleteNode(*manager.NodeInstance) error
}

type instanceInsertQuery func(storage.Transaction, *instancedb.Record) error

type instanceSelectByUserIDQuery func(storage.Transaction, *instancedb.SelectByUserIDQueryInput) (*instancedb.Record, error)

type instanceSelectAllByUserIDQuery func(storage.Transaction, *instancedb.SelectAllByUserIDQuery) ([]*instancedb.Record, error)

type instanceUpdateQuery func(storage.Transaction, *instancedb.UpdateQueryInput) error

type instanceSelectQuery func(storage.Transaction, *instancedb.SelectQueryInput) (*instancedb.Record, error)
