package engine

import (
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
	"github.com/stretchr/testify/mock"
)

type DataSource interface {
	FetchFacts() (map[string]interface{}, error)
}

type MockDataSource struct {
	mock.Mock
}

func (m *MockDataSource) FetchFacts() (map[string]interface{}, error) {
	args := m.Called()
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

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

// TestEvaluateWithReportFacts is a unit test for the EvaluateWithReportFacts function.
//
// It creates a new engine and enables the ReportFacts option.
// Then it defines a rule with a condition and an event.
// The rule is added to the engine and a fact is defined.
// The fact is evaluated using the engine, and the resulting events are checked.
// Finally, the function asserts that the event includes the triggering fact and its value.
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

func TestIntegrationEngineWithRealWorldScenario(t *testing.T) {
	engine := NewEngine()

	// Add rules
	rule1 := rules.Rule{
		Name:     "Rule 1",
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
			EventType: "High Temperature",
		},
	}
	err := engine.AddRule(rule1)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	rule2 := rules.Rule{
		Name:     "Rule 2",
		Priority: 2,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "humidity",
					Operator: "greaterThan",
					Value:    70,
				},
			},
		},
		Event: rules.Event{
			EventType: "High Humidity",
		},
	}
	err = engine.AddRule(rule2)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Evaluate facts
	fact := rules.Fact{
		"temperature": 35,
		"humidity":    75,
	}
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	// Update rule
	rule1Updated := rule1
	rule1Updated.Priority = 3
	err = engine.UpdateRule("Rule 1", rule1Updated)
	if err != nil {
		t.Fatalf("Failed to update rule: %v", err)
	}

	// Remove rule
	err = engine.RemoveRule("Rule 2")
	if err != nil {
		t.Fatalf("Failed to remove rule: %v", err)
	}

	// Evaluate facts again
	events, err = engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

}

func TestIntegrationEngineWithExternalData(t *testing.T) {
	// Mock external data source
	mockDataSource := new(MockDataSource)

	// Set up the mock data source to return specific facts when called
	mockDataSource.On("FetchFacts").Return(map[string]interface{}{
		"temperature": 35,
		"humidity":    75,
	}, nil)

	// Create a new engine
	engine := NewEngine()

	// Add rules
	rule1 := rules.Rule{
		Name:     "Rule 1",
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
			EventType: "High Temperature",
		},
	}
	err := engine.AddRule(rule1)
	if err != nil {
		t.Fatalf("Failed to add rule: %v", err)
	}

	// Fetch facts from the mock data source
	facts, err := mockDataSource.FetchFacts()
	if err != nil {
		t.Fatalf("Failed to fetch facts: %v", err)
	}

	// Evaluate facts
	events, err := engine.Evaluate(facts)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	// Verify the results
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "High Temperature" {
		t.Fatalf("Expected event type 'High Temperature', got '%s'", events[0].EventType)
	}
}

func TestEngine_EvaluateRules_InvalidRule(t *testing.T) {
	engine := NewEngine()

	invalidRule := rules.Rule{
		Name:     "InvalidRule",
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
			EventType: "High Temperature",
		},
	}

	// Directly add the invalid rule to the engine's Rules map
	engine.Rules[invalidRule.Name] = invalidRule

	// Also add the invalid rule to the engine's RuleIndex map
	engine.RuleIndex["temperature"] = append(engine.RuleIndex["temperature"], &invalidRule)

	fact := rules.Fact{
		"temperature": 35,
	}

	_, err := engine.Evaluate(fact)
	if err == nil {
		t.Fatalf("Expected error when evaluating invalid rule, but got none")
	}
}

func TestEngine_EvaluateRules_MixedValidity(t *testing.T) {
	engine := NewEngine()

	// Add a valid rule
	validRule := rules.Rule{
		Name:     "Valid Rule",
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
			EventType: "High Temperature",
		},
	}
	err := engine.AddRule(validRule)
	if err != nil {
		t.Fatalf("Failed to add valid rule: %v", err)
	}

	// Add an invalid rule
	invalidRule := rules.Rule{
		Name:     "Invalid Rule",
		Priority: 2,
		Conditions: rules.Conditions{
			All: []rules.Condition{
				{
					Fact:     "humidity",
					Operator: "invalidOperator", // This operator is invalid
					Value:    70,
				},
			},
		},
		Event: rules.Event{
			EventType: "High Humidity",
		},
	}
	err = engine.AddRule(invalidRule)
	if err == nil {
		t.Fatalf("Expected error when adding invalid rule, got nil")
	}

	// Evaluate facts
	fact := rules.Fact{
		"temperature": 35,
		"humidity":    75,
	}
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	// Only the valid rule should have been evaluated
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "High Temperature" {
		t.Fatalf("Expected event type 'High Temperature', got '%s'", events[0].EventType)
	}
}

func TestEvaluateWithNestedConditions(t *testing.T) {
	engine := NewEngine()

	// Define your nested conditions
	rule := rules.Rule{
		Name:     "Nested Rule",
		Priority: 1,
		Conditions: rules.Conditions{
			Any: []rules.Condition{
				{
					All: []rules.Condition{
						{
							Fact:     "temperature",
							Operator: "greaterThan",
							Value:    30,
						},
						{
							Fact:     "humidity",
							Operator: "lessThan",
							Value:    70,
						},
					},
				},
				{
					All: []rules.Condition{
						{
							Fact:     "windSpeed",
							Operator: "greaterThan",
							Value:    10,
						},
						{
							Fact:     "rainfall",
							Operator: "greaterThan",
							Value:    20,
						},
					},
				},
			},
		},
		Event: rules.Event{
			EventType: "Complex Weather Condition",
		},
	}

	// Add the rule to the engine
	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add nested rule: %v", err)
	}

	// Define the facts that will be used for evaluation
	fact := rules.Fact{
		"temperature": 35,
		"humidity":    65,
		"windSpeed":   15,
		"rainfall":    25,
	}

	// Evaluate the rules with the facts
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	// Check the result
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "Complex Weather Condition" {
		t.Fatalf("Expected event type 'Complex Weather Condition', got '%s'", events[0].EventType)
	}
}

func TestEvaluateWithDeeplyNestedConditions(t *testing.T) {
	engine := NewEngine()

	// Define your deeply nested conditions
	rule := rules.Rule{
		Name:     "Deeply Nested Rule",
		Priority: 1,
		Conditions: rules.Conditions{
			Any: []rules.Condition{
				{
					All: []rules.Condition{
						{
							Fact:     "temperature",
							Operator: "greaterThan",
							Value:    30,
						},
						{
							Any: []rules.Condition{
								{
									Fact:     "humidity",
									Operator: "lessThan",
									Value:    70,
								},
								{
									Fact:     "pressure",
									Operator: "greaterThan",
									Value:    1000,
								},
							},
						},
					},
				},
				{
					All: []rules.Condition{
						{
							Fact:     "windSpeed",
							Operator: "greaterThan",
							Value:    10,
						},
						{
							All: []rules.Condition{
								{
									Fact:     "rainfall",
									Operator: "greaterThan",
									Value:    20,
								},
								{
									Fact:     "cloudCover",
									Operator: "greaterThan",
									Value:    50,
								},
							},
						},
					},
				},
			},
		},
		Event: rules.Event{
			EventType: "Complex Weather Condition",
		},
	}

	// Add the rule to the engine
	err := engine.AddRule(rule)
	if err != nil {
		t.Fatalf("Failed to add deeply nested rule: %v", err)
	}

	// Define the facts that will be used for evaluation
	fact := rules.Fact{
		"temperature": 35,
		"humidity":    65,
		"pressure":    1010,
		"windSpeed":   15,
		"rainfall":    25,
		"cloudCover":  60,
	}

	// Evaluate the rules with the facts
	events, err := engine.Evaluate(fact)
	if err != nil {
		t.Fatalf("Failed to evaluate facts: %v", err)
	}

	// Check the result
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "Complex Weather Condition" {
		t.Fatalf("Expected event type 'Complex Weather Condition', got '%s'", events[0].EventType)
	}
}
