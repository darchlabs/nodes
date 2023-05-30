package api

import (
	"github.com/darchlabs/nodes/internal/manager"
	"github.com/darchlabs/nodes/internal/storage"
	"github.com/darchlabs/nodes/internal/storage/instance"
)

type mockNodeManager struct {
	res *manager.NodeInstance
	err error
}

func (m *mockNodeManager) DeployNewNode(_ *manager.CreateDeploymentOptions) (*manager.NodeInstance, error) {
	return m.res, m.err
}

func (m *mockNodeManager) DeleteArtifacts(_ *manager.Artifacts) error {
	return m.err
}

type mockInstanceInsertQuery struct {
	err error
}

func (m *mockInstanceInsertQuery) Insert(_ storage.Transaction, _ *instance.Record) error {
	return m.err
}
