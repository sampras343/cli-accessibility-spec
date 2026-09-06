package domains

import (
	"context"
	"testing"

	"github.com/sampras343/cli-accessibility-spec/internal/engine"
)

type mockDomain struct {
	name     string
	criteria []Criterion
}

func (d *mockDomain) Name() string       { return d.name }
func (d *mockDomain) Criteria() []Criterion { return d.criteria }

type mockCriterion struct {
	id     string
	domain string
}

func (c *mockCriterion) ID() string                { return c.id }
func (c *mockCriterion) Name() string              { return c.id + " test" }
func (c *mockCriterion) Domain() string            { return c.domain }
func (c *mockCriterion) Level() engine.Level       { return engine.LevelA }
func (c *mockCriterion) Testability() engine.Testability { return engine.Auto }
func (c *mockCriterion) SpecVersion() string       { return "1.0" }
func (c *mockCriterion) Precondition(_ *ProbeResult) bool { return true }
func (c *mockCriterion) Run(_ context.Context, _ string, _ *ProbeResult) *engine.Result {
	return &engine.Result{ID: c.id, Outcome: engine.Supports}
}

func TestRegisterAndRetrieve(t *testing.T) {
	resetRegistry()
	Register("test", &mockDomain{
		name:     "test",
		criteria: []Criterion{&mockCriterion{id: "T-1", domain: "test"}},
	})
	domains := AllDomains()
	if len(domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(domains))
	}
	if domains[0].Name() != "test" {
		t.Errorf("domain name = %q, want %q", domains[0].Name(), "test")
	}
	if len(domains[0].Criteria()) != 1 {
		t.Fatalf("expected 1 criterion, got %d", len(domains[0].Criteria()))
	}
	if domains[0].Criteria()[0].ID() != "T-1" {
		t.Errorf("criterion ID = %q, want %q", domains[0].Criteria()[0].ID(), "T-1")
	}
}

func TestRegisterPanicsOnDuplicate(t *testing.T) {
	resetRegistry()
	Register("dup", &mockDomain{name: "dup"})
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate registration")
		}
	}()
	Register("dup", &mockDomain{name: "dup"})
}
