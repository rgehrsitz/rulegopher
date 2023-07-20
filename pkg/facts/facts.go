package facts

import (
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

// FactHandler is used to handle facts and evaluate them using an engine.
type FactHandler struct {
	engine *engine.Engine
}

// NewFactHandler creates a new FactHandler with the given engine.
func NewFactHandler(engine *engine.Engine) *FactHandler {
	return &FactHandler{
		engine: engine,
	}
}

// HandleFact evaluates the given fact using the FactHandler's engine.
func (factHandler *FactHandler) HandleFact(fact rules.Fact) ([]rules.Event, error) {
	return factHandler.engine.Evaluate(fact)
}
