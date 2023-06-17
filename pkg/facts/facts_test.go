package facts

import (
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func TestHandleFact(t *testing.T) {
	e := engine.NewEngine()
	fh := NewFactHandler(e)

	// Define a rule
	rule := rules.Rule{
		Name:     "TestRule",
		Priority: 1,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "temperature",
					Operator: "greaterThan",
					Value:    30,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	// Add the rule to the engine
	e.AddRule(rule)

	// Define a fact that satisfies the rule
	fact := rules.Fact{
		"temperature": 35,
	}

	// Handle the fact
	events := fh.HandleFact(fact)

	// Check that the correct event was generated
	if len(events) != 1 || events[0].EventType != "alert" {
		t.Errorf("HandleFact returned incorrect events: got %v", events)
	}
}
