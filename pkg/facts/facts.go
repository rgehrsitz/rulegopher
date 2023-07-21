package facts

import (
	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

// The FactHandler type is a struct that contains a reference to an engine.Engine object.
// @property engine - The `engine` property is a pointer to an instance of the `Engine` struct.
type FactHandler struct {
	engine *engine.Engine
}

// The NewFactHandler function returns a new instance of the FactHandler struct with the provided
// engine.
func NewFactHandler(engine *engine.Engine) *FactHandler {
	return &FactHandler{
		engine: engine,
	}
}

// The `HandleFact` function is a method of the `FactHandler` struct. It takes a `fact` of type
// `rules.Fact` as a parameter and returns a slice of `rules.Event` and an error.
func (factHandler *FactHandler) HandleFact(fact rules.Fact) ([]rules.Event, error) {
	return factHandler.engine.Evaluate(fact)
}
