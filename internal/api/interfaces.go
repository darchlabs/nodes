package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/internal/storage/instance"
)

type nodeManager interface {
	DeployNewNode(*manager.CreateDeploymentOptions) (*manager.NodeInstance, error)
	DeleteNode(*manager.NodeInstance) error
}

type instanceInsertQuery func(storage.Transaction, *instance.Record) error

type instanceSelectByUserIDQuery func(storage.Transaction, *instance.SelectByUserIDQueryInput) (*instance.Record, error)

type instanceSelectAllQuery func(storage.Transaction) ([]*instance.Record, error)

type instanceUpdateQuery func(storage.Transaction, *instance.UpdateQueryInput) error
