package rules

import (
	"testing"
)

func TestConditionEvaluate(t *testing.T) {
	// Define a condition
	condition := Condition{
		Fact:     "temperature",
		Operator: "greaterThan",
		Value:    30,
	}

	// Define a fact where the condition should be true
	factTrue := Fact{
		"temperature": 40,
	}

	// Define a fact where the condition should be false
	factFalse := Fact{
		"temperature": 20,
	}

	// Test the condition with the fact where it should be true
	satisfied, _, _ := condition.Evaluate(factTrue)
	if !satisfied {
		t.Errorf("Expected condition to be true, but it was false")
	}

	// Test the condition with the fact where it should be false
	satisfied, _, _ = condition.Evaluate(factFalse)
	if satisfied {
		t.Errorf("Expected condition to be false, but it was true")
	}
}

func TestConditionEvaluate2(t *testing.T) {
	// Define conditions
	conditions := []Condition{
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
		{
			Fact:     "location",
			Operator: "equal",
			Value:    "indoors",
		},
		{
			Fact:     "motionDetected",
			Operator: "notEqual",
			Value:    false,
		},
	}

	// Define a fact where all conditions should be true
	factTrue := Fact{
		"temperature":    40,
		"humidity":       0.4,
		"location":       "indoors",
		"motionDetected": true,
	}

	// Define a fact where all conditions should be false
	factFalse := Fact{
		"temperature":    20,
		"humidity":       0.6,
		"location":       "outdoors",
		"motionDetected": false,
	}

	// Test each condition with the fact where it should be true
	for _, condition := range conditions {
		satisfied, _, _ := condition.Evaluate(factTrue)
		if !satisfied {
			t.Errorf("Expected condition to be true, but it was false")
		}
	}

	// Test each condition with the fact where it should be false
	for _, condition := range conditions {
		satisfied, _, _ := condition.Evaluate(factFalse)
		if satisfied {
			t.Errorf("Expected condition to be false, but it was true")
		}
	}
}

func TestRuleEvaluate3(t *testing.T) {
	// Define a rule with nested 'any' and 'all' conditions
	rule := Rule{
		Name:     "TestRule",
		Priority: 1,
		Conditions: Conditions{
			All: []Condition{
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
			Any: []Condition{
				{
					Fact:     "location",
					Operator: "equal",
					Value:    "indoors",
				},
				{
					Fact:     "motionDetected",
					Operator: "equal",
					Value:    true,
				},
			},
		},
		Event: Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	// Define a fact where the rule should be satisfied
	factTrue := Fact{
		"temperature":    40,
		"humidity":       0.4,
		"location":       "indoors",
		"motionDetected": false,
	}

	// Define a fact where the rule should not be satisfied
	factFalse := Fact{
		"temperature":    20,
		"humidity":       0.6,
		"location":       "outdoors",
		"motionDetected": false,
	}

	// Test the rule with the fact where it should be satisfied
	if !rule.Evaluate(factTrue, true) {
		t.Errorf("Expected rule to be satisfied, but it was not")
	}

	// Test the rule with the fact where it should not be satisfied
	if rule.Evaluate(factFalse, true) {
		t.Errorf("Expected rule to not be satisfied, but it was")
	}
}

func TestRuleEvaluateComplex(t *testing.T) {
	// Define a complex rule with nested 'any' and 'all' conditions
	rule := Rule{
		Name:     "TestRuleComplex",
		Priority: 1,
		Conditions: Conditions{
			All: []Condition{
				{
					Fact:     "temperature",
					Operator: "greaterThan",
					Value:    30,
					All: []Condition{
						{
							Fact:     "humidity",
							Operator: "lessThan",
							Value:    0.5,
						},
					},
				},
			},
			Any: []Condition{
				{
					Fact:     "location",
					Operator: "equal",
					Value:    "indoors",
					Any: []Condition{
						{
							Fact:     "motionDetected",
							Operator: "equal",
							Value:    true,
						},
					},
				},
			},
		},
		Event: Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	// Define a fact where the rule should be satisfied
	factTrue := Fact{
		"temperature":    35,
		"humidity":       0.4,
		"location":       "indoors",
		"motionDetected": true,
	}

	// Define a fact where the rule should not be satisfied
	factFalse := Fact{
		"temperature":    20,
		"humidity":       0.6,
		"location":       "outdoors",
		"motionDetected": false,
	}

	// Test the rule with the fact where it should be satisfied
	if !rule.Evaluate(factTrue, true) {
		t.Errorf("Expected rule to be satisfied, but it was not")
	}

	// Test the rule with the fact where it should not be satisfied
	if rule.Evaluate(factFalse, true) {
		t.Errorf("Expected rule to not be satisfied, but it was")
	}
}

func TestConditionEvaluateWithDifferentOperators(t *testing.T) {
	// Define a fact
	fact := Fact{
		"temperature": 35,
		"humidity":    0.4,
		"location":    "indoors",
		"motion":      []string{"run", "walk"},
	}

	// Define conditions with different operators
	conditions := []Condition{
		{
			Fact:     "temperature",
			Operator: "equal",
			Value:    35,
		},
		{
			Fact:     "temperature",
			Operator: "notEqual",
			Value:    30,
		},
		{
			Fact:     "temperature",
			Operator: "greaterThan",
			Value:    30,
		},
		{
			Fact:     "temperature",
			Operator: "greaterThanOrEqual",
			Value:    35,
		},
		{
			Fact:     "humidity",
			Operator: "lessThan",
			Value:    0.5,
		},
		{
			Fact:     "humidity",
			Operator: "lessThanOrEqual",
			Value:    0.4,
		},
		{
			Fact:     "location",
			Operator: "contains",
			Value:    "door",
		},
		{
			Fact:     "motion",
			Operator: "notContains",
			Value:    "jump",
		},
	}

	// Test each condition with the fact
	for _, condition := range conditions {
		satisfied, _, _ := condition.Evaluate(fact)
		if !satisfied {
			t.Errorf("Expected condition to be true, but it was false")
		}
	}
}

func TestConditionEvaluateWithDifferentValueTypes(t *testing.T) {
	// Define a fact
	fact := Fact{
		"temperature": 35.5,
		"humidity":    0.4,
		"location":    "indoors",
		"people":      5,
	}

	// Define conditions with different types of values
	conditions := []Condition{
		{
			Fact:     "temperature",
			Operator: "greaterThan",
			Value:    30.0,
		},
		{
			Fact:     "humidity",
			Operator: "lessThan",
			Value:    0.5,
		},
		{
			Fact:     "location",
			Operator: "equal",
			Value:    "indoors",
		},
		{
			Fact:     "people",
			Operator: "greaterThan",
			Value:    3,
		},
	}

	// Test each condition with the fact
	for _, condition := range conditions {
		satisfied, _, _ := condition.Evaluate(fact)
		if !satisfied {
			t.Errorf("Expected condition to be true, but it was false")
		}
	}
}

func TestConditionEvaluateWithMissingFact(t *testing.T) {
	// Define a fact
	fact := Fact{
		"temperature": 35.5,
		"humidity":    0.4,
		"location":    "indoors",
	}

	// Define a condition that requires a fact not present in the fact map
	condition := Condition{
		Fact:     "people",
		Operator: "greaterThan",
		Value:    3,
	}

	// Test the condition with the fact
	satisfied, _, _ := condition.Evaluate(fact)
	if satisfied {
		t.Errorf("Expected condition to be false, but it was true")
	}
}
