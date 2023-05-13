package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/internal/storage/instance"
)

type nodeManager interface {
	DeployNewNode(*manager.CreateDeploymentOptions) (*manager.NodeInstance, error)
}

type instanceInsertQuery func(storage.Transaction, *instance.Record) error

type instanceSelectQuery func(storage.Transaction, *instance.SelectQueryInput) (*instance.Record, error)
