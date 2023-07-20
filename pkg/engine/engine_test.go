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
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "temperature",
					Operator: "greaterThan",
					Value:    30,
				},
			},
		},
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
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

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
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that the list of events contains both events
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	// Check that the events are ordered by rule priority
	if events[0].EventType != "TestEvent2" || events[1].EventType != "TestEvent1" {
		t.Fatalf("Events are not ordered by rule priority")
	}
}

func TestExecuteEmptyRules(t *testing.T) {
	// Create a new engine
	e := NewEngine()

	// Define a fact
	fact := rules.Fact{
		"temperature": 35,
	}

	// Execute rules with the fact
	events, err := e.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that no events were returned
	if len(events) != 0 {
		t.Errorf("ExecuteRules returned events when no rules were added: got %v", events)
	}
}

func TestAddRuleWithNilCondition(t *testing.T) {
	engine := NewEngine()

	rule := rules.Rule{
		Name:     "TestRule",
		Priority: 1,
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	// Try to add the rule with a nil condition
	err := engine.AddRule(rule)
	if err == nil {
		t.Fatalf("Expected error when adding rule with nil condition, got nil")
	}
}

func TestEvaluateWithEmptyFact(t *testing.T) {
	engine := NewEngine()

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

	// Add the rule
	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Define an empty fact
	fact := rules.Fact{}

	// Evaluate the fact
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that the list of events is empty
	if len(events) != 0 {
		t.Fatalf("Expected no events, got %d", len(events))
	}
}

func TestEvaluateWithMultipleMatchingFacts(t *testing.T) {
	engine := NewEngine()

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
				{
					Fact:     "humidity",
					Operator: "lessThan",
					Value:    0.5,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	fact := rules.Fact{
		"temperature": 35,
		"humidity":    0.4,
	}

	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].EventType != "alert" {
		t.Errorf("Expected event type 'alert', got '%s'", events[0].EventType)
	}
}

func TestEvaluateWithNoMatchingFacts(t *testing.T) {
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
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Define a fact that does not match the rule
	fact := rules.Fact{
		"temperature": 20,
	}

	// Evaluate the fact
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that no events were generated
	if len(events) != 0 {
		t.Errorf("Expected no events, got %d", len(events))
	}
}

func TestUpdateRuleWithDifferentPriority(t *testing.T) {
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

	// Add the rule
	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Update the rule with a different priority
	newRule := rule
	newRule.Priority = 2
	err = engine.UpdateRule(rule.Name, newRule)
	if err != nil {
		t.Fatalf("Failed to update rule: %v", err)
	}

	// Check that the rule's priority has been updated
	for _, r := range engine.Rules {
		if r.Name == rule.Name {
			if r.Priority != newRule.Priority {
				t.Errorf("Expected rule priority to be %d, got %d", newRule.Priority, r.Priority)
			}
			break
		}
	}
}

func TestRemoveRuleWithMultipleRules(t *testing.T) {
	engine := NewEngine()

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
		Priority: 2,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "humidity",
					Operator: "lessThan",
					Value:    0.5,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "Dehumidifier turned on",
		},
	}

	// Add the rules
	err := engine.AddRule(rule1)
	if err != nil {
		t.Fatalf("Failed to add rule1: %v", err)
	}

	err = engine.AddRule(rule2)
	if err != nil {
		t.Fatalf("Failed to add rule2: %v", err)
	}

	// Remove rule1
	err = engine.RemoveRule(rule1.Name)
	if err != nil {
		t.Fatalf("Failed to remove rule1: %v", err)
	}

	// Define a fact that matches rule2's conditions
	fact := rules.Fact{
		"humidity": 0.4,
	}

	// Evaluate the fact
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that the list of events contains only rule2's event
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].EventType != rule2.Event.EventType {
		t.Fatalf("Expected event from rule2, got event from another rule")
	}
}

func TestEvaluateWithReportFacts(t *testing.T) {
	// Create a new engine
	engine := NewEngine()
	engine.ReportFacts = true // Enable the ReportFacts option

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
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Define a fact that should trigger the rule
	fact := rules.Fact{
		"temperature": 35,
	}

	// Evaluate the fact
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	// Check that the event includes the triggering fact
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	event := events[0]
	if len(event.Facts) != 1 || event.Facts[0] != "temperature" {
		t.Errorf("Event does not include the triggering fact")
	}
	if len(event.Values) != 1 || event.Values[0] != 35 {
		t.Errorf("Event does not include the triggering fact value")
	}
}

func TestEvaluateWithReportRuleName(t *testing.T) {
	engine := NewEngine()
	engine.ReportRuleName = true

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

	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	fact := rules.Fact{
		"temperature": 35,
	}

	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Error evaluating fact: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].RuleName != "TestRule" {
		t.Fatalf("Expected RuleName to be 'TestRule', got '%s'", events[0].RuleName)
	}
}

func TestAddRuleWithInvalidOperator(t *testing.T) {
	engine := NewEngine()

	rule := rules.Rule{
		Name:     "TestRule",
		Priority: 1,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "temperature",
					Operator: "invalidOperator",
					Value:    30,
				},
			},
		},
		Event: rules.Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	err := engine.AddRule(rule)
	if err == nil {
		t.Fatalf("Expected error when adding rule with invalid operator, got nil")
	}
}

func TestUpdateRuleWithInvalidName(t *testing.T) {
	engine := NewEngine()

	rule := rules.Rule{
		Name: "TestRule",
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

	// Try to update a rule that doesn't exist
	err := engine.UpdateRule("NonExistentRule", rule)
	if err == nil {
		t.Fatalf("Expected error when updating non-existent rule, got nil")
	}
}

func TestRemoveRuleWithInvalidName(t *testing.T) {
	engine := NewEngine()

	// Try to remove a rule that doesn't exist
	err := engine.RemoveRule("NonExistentRule")
	if err == nil {
		t.Fatalf("Expected error when removing non-existent rule, got nil")
	}
}

func TestAddRuleWithEmptyName(t *testing.T) {
	engine := NewEngine()

	rule := rules.Rule{
		Name:     "",
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
			EventType: "alert",
		},
	}

	err := engine.AddRule(rule)
	if err == nil {
		t.Errorf("Expected error when adding rule with empty name, but got none")
	}
}
