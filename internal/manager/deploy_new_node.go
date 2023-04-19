package manager

import "github.com/pkg/errors"

type CreatePodOptions struct {
	Network string
	EnvVars map[string]string
}

func (m *Manager) DeployNewNode(opts *CreatePodOptions) (*NodeInstance, error) {
	setupFunc, ok := m.networkNodeSetups[opts.Network]
	if !ok {
		return nil, ErrNodeNotFound
	}

	nodeInstace, err := setupFunc(opts.Network, opts.EnvVars)
	if err != nil {
		return nil, errors.Wrapf(err, "manager: Manager.DeployNewNode for %s network", opts.Network)
	}

	return nodeInstace, nil
}

func intPtr32(i int32) *int32 {
	return &i
}
