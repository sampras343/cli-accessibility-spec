package domains

import (
	"context"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
	"github.com/sampras343/cli-accessibility-spec/internal/probe"
)

type Criterion interface {
	ID() string
	Name() string
	Domain() string
	Level() engine.Level
	Testability() engine.Testability
	SpecVersion() string
	Precondition(probe *probe.ProbeResult) bool
	Run(ctx context.Context, binary string, probe *probe.ProbeResult) *engine.Result
}

type Domain interface {
	Name() string
	Criteria() []Criterion
}
