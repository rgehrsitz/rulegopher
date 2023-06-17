package engine

import (
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func TestAddRule(t *testing.T) {
	// Create a new engine
	engine := NewEngine()

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
	err := engine.AddRule(rule)
	if err != nil {
		t.Errorf("Failed to add rule: %v", err)
	}

	// Check whether the rule exists in the engine's rule set
	if _, exists := engine.rules[rule.Name]; !exists {
		t.Errorf("Rule was not added to the engine")
	}
}
