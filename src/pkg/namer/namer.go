package namer

import (
	"math/rand"

	"github.com/lucasepe/codename"
	"github.com/pkg/errors"
)

type Namer struct {
	rng *rand.Rand
}

func New() (*Namer, error) {
	rng, err := codename.DefaultRNG()
	if err != nil {
		return nil, errors.Wrap(err, "namer: New codename.DefaultRNG error")
	}
	return &Namer{
		rng: rng,
	}, nil
}

func (n *Namer) Generate() string {
	return codename.Generate(n.rng, 0)
}
