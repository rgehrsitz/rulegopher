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

func TestFactHandlerHandleFactWithNoSatisfyingRules(t *testing.T) {
	// Create a new engine
	eng := engine.NewEngine()

	// Create a new fact handler with the engine
	fh := NewFactHandler(eng)

	// Define a fact that does not satisfy any rules
	fact := rules.Fact{
		"temperature": 20,
		"humidity":    0.6,
		"location":    "outdoors",
	}

	// Handle the fact
	events := fh.HandleFact(fact)

	// Check that no events were returned
	if len(events) != 0 {
		t.Errorf("Expected no events to be returned, but got %d", len(events))
	}
}

func TestFactHandlerHandleFactWithMultipleSatisfyingRules(t *testing.T) {
	// Create a new engine
	eng := engine.NewEngine()

	// Add some rules to the engine
	rule1 := rules.Rule{
		Name:     "TestRule1",
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
	rule2 := rules.Rule{
		Name:     "TestRule2",
		Priority: 1,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "humidity",
					Operator: "greaterThan",
					Value:    0.5,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "Dehumidifier turned on",
		},
	}
	eng.AddRule(rule1)
	eng.AddRule(rule2)

	// Create a new fact handler with the engine
	fh := NewFactHandler(eng)

	// Define a fact that satisfies both rules
	fact := rules.Fact{
		"temperature": 35,
		"humidity":    0.6,
	}

	// Handle the fact
	events := fh.HandleFact(fact)

	// Check that two events were returned
	if len(events) != 2 {
		t.Errorf("Expected two events to be returned, but got %d", len(events))
	}
}
