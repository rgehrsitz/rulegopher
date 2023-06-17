package engine

import (
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func TestEngine(t *testing.T) {
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

	// Test AddRule
	err := engine.AddRule(rule)
	if err != nil {
		t.Errorf("Failed to add rule: %v", err)
	}

	// Test AddRule with an existing rule
	err = engine.AddRule(rule)
	if err == nil {
		t.Errorf("Expected error when adding an existing rule, got nil")
	}

	// Test RemoveRule
	err = engine.RemoveRule(rule.Name)
	if err != nil {
		t.Errorf("Failed to remove rule: %v", err)
	}

	// Test RemoveRule with a non-existing rule
	err = engine.RemoveRule(rule.Name)
	if err == nil {
		t.Errorf("Expected error when removing a non-existing rule, got nil")
	}

	// Test UpdateRule with a non-existing rule
	err = engine.UpdateRule(rule.Name, rule)
	if err == nil {
		t.Errorf("Expected error when updating a non-existing rule, got nil")
	}

	// Add the rule back to test UpdateRule with an existing rule
	engine.AddRule(rule)
	err = engine.UpdateRule(rule.Name, rule)
	if err != nil {
		t.Errorf("Failed to update rule: %v", err)
	}
}
