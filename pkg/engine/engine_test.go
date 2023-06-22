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

func TestAddDuplicateRule(t *testing.T) {
	engine := NewEngine()

	rule := rules.Rule{
		Name: "TestRule",
		// ... other properties ...
	}

	// Add the rule for the first time
	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Try to add the same rule again
	err = engine.AddRule(rule)
	if err == nil {
		t.Fatalf("Expected error when adding duplicate rule, got nil")
	}
}

func TestRemoveNonExistentRule(t *testing.T) {
	engine := NewEngine()

	// Try to remove a rule that doesn't exist
	err := engine.RemoveRule("NonExistentRule")
	if err == nil {
		t.Fatalf("Expected error when removing non-existent rule, got nil")
	}
}

func TestUpdateNonExistentRule(t *testing.T) {
	engine := NewEngine()

	newRule := rules.Rule{
		Name: "NewRule",
		// ... other properties ...
	}

	// Try to update a rule that doesn't exist
	err := engine.UpdateRule("NonExistentRule", newRule)
	if err == nil {
		t.Fatalf("Expected error when updating non-existent rule, got nil")
	}
}

func TestEvaluateNonMatchingFact(t *testing.T) {
	engine := NewEngine()

	fact := rules.Fact{
		"NonMatchingFact": "value",
		// ... other properties ...
	}

	// Evaluate the fact
	events := engine.Evaluate(fact)

	// Check that the list of events is empty
	if len(events) != 0 {
		t.Fatalf("Expected no events, got %d", len(events))
	}
}

func TestEvaluateFactMatchingMultipleRules(t *testing.T) {
	engine := NewEngine()

	rule1 := rules.Rule{
		Name:     "TestRule1",
		Priority: 2,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "TestFact",
					Operator: "equal",
					Value:    "value",
				},
			},
		},
		Event: rules.Event{
			EventType: "TestEvent1",
		},
	}

	rule2 := rules.Rule{
		Name:     "TestRule2",
		Priority: 1,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "TestFact",
					Operator: "equal",
					Value:    "value",
				},
			},
		},
		Event: rules.Event{
			EventType: "TestEvent2",
		},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)

	fact := rules.Fact{
		"TestFact": "value",
	}

	// Evaluate the fact
	events := engine.Evaluate(fact)

	// Check that the list of events contains both events
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	// Check that the events are ordered by rule priority
	if events[0].EventType != "TestEvent2" || events[1].EventType != "TestEvent1" {
		t.Fatalf("Events are not ordered by rule priority")
	}
}
