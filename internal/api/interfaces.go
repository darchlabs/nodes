package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/internal/storage/instance"
)

type nodeManager interface {
	DeployNewNode(*manager.CreateDeploymentOptions) (*manager.NodeInstance, error)
}

type instanceInsertQuery func(tx storage.Transaction, record *instance.Record) error
