package facts

import (
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type FactHandler struct {
	engine *engine.Engine
}

func NewFactHandler(engine *engine.Engine) *FactHandler {
	return &FactHandler{
		engine: engine,
	}
}

func (fh *FactHandler) HandleFact(fact rules.Fact) []rules.Event {
	return fh.engine.Evaluate(fact)
}
