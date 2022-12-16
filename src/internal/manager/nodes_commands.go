package manager

import (
	"fmt"

	"github.com/darchlabs/nodes/src/internal/command"
)

type nodeRunner func(*NodeConfig) *command.Command

var (
	NetworkEthereum = "ethereum"
	NetworkPolygon  = "polygon"

	networkNodesRunners = map[string]nodeRunner{
		NetworkEthereum: newGanacheCommand,
		NetworkPolygon:  newGanacheCommand,
	}
)

func newGanacheCommand(config *NodeConfig) *command.Command {
	bootstrapNodeURL := config.BootsrapNodeURL
	if config.FromBlockNumber != 0 {
		bootstrapNodeURL = fmt.Sprintf("%s@%d", config.BootsrapNodeURL, config.FromBlockNumber)
	}

	return command.New(
		"ganache",
		"-p", fmt.Sprintf("%d", config.Port),
		"--host", config.Host,
		"--db", config.BaseChainDataPath,
		"--fork", bootstrapNodeURL,
	)
}
