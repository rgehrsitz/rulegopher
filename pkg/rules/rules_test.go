package rules

import (
	"testing"
)

// TestConditionEvaluate tests the Evaluate function of the Condition struct.
//
// It defines a condition with a specific fact, operator, and value.
// It also defines two facts, one where the condition should be true and one where it should be false.
// The function tests the condition with both facts and verifies the expected results.
// It reports any error encountered during the evaluation.
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
	satisfied, _, _, err := condition.Evaluate(factTrue, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected condition to be true, but it was false")
	}

	// Test the condition with the fact where it should be false
	satisfied, _, _, err = condition.Evaluate(factFalse, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
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
		satisfied, _, _, err := condition.Evaluate(factTrue, "Ignore")
		if err != nil {
			t.Fatalf("Error evaluating condition: %v", err)
		}
		if !satisfied {
			t.Errorf("Expected condition to be true, but it was false")
		}
	}

	// Test each condition with the fact where it should be false
	for _, condition := range conditions {
		satisfied, _, _, err := condition.Evaluate(factFalse, "Ignore")
		if err != nil {
			t.Fatalf("Error evaluating condition: %v", err)
		}
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
	satisfied, err := rule.Evaluate(factTrue, true, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating rule: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected rule to be satisfied, but it was not")
	}

	// Test the rule with the fact where it should not be satisfied

	satisfied, err = rule.Evaluate(factFalse, true, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating rule: %v", err)
	}
	if satisfied {
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
	satisfied, err := rule.Evaluate(factTrue, true, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating rule: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected rule to be satisfied, but it was not")
	}

	// Test the rule with the fact where it should not be satisfied

	satisfied, err = rule.Evaluate(factFalse, true, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating rule: %v", err)
	}
	if satisfied {
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
		satisfied, _, _, err := condition.Evaluate(fact, "Ignore")
		if err != nil {
			t.Fatalf("Error evaluating condition: %v", err)
		}
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
		satisfied, _, _, err := condition.Evaluate(fact, "Ignore")
		if err != nil {
			t.Fatalf("Error evaluating condition: %v", err)
		}
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
	satisfied, _, _, err := condition.Evaluate(fact, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if satisfied {
		t.Errorf("Expected condition to be false, but it was true")
	}
}

func TestEvaluateWithInvalidFact(t *testing.T) {
	// Define a rule
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
			},
		},
		Event: Event{
			EventType:      "alert",
			CustomProperty: "AC turned on",
		},
	}

	// Define a fact with invalid data type for the condition
	fact := Fact{
		"temperature": "forty", // string instead of a number
	}

	// Test the rule with the fact
	_, err := rule.Evaluate(fact, true, "Ignore")
	if err == nil {
		t.Errorf("Expected an error due to invalid fact data type, but got none")
	}
}

func TestConditionEvaluateDeeplyNested(t *testing.T) {
	// Define a condition with deeply nested conditions
	condition := Condition{
		All: []Condition{
			{
				Fact:     "temperature",
				Operator: "greaterThan",
				Value:    30,
				Any: []Condition{
					{
						Fact:     "humidity",
						Operator: "lessThan",
						Value:    0.5,
						All: []Condition{
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
				},
			},
		},
	}

	// Define a fact where the condition should be true
	factTrue := Fact{
		"temperature":    40,
		"humidity":       0.4,
		"location":       "indoors",
		"motionDetected": true,
	}

	// Define a fact where the condition should be false
	factFalse := Fact{
		"temperature":    20,
		"humidity":       0.6,
		"location":       "outdoors",
		"motionDetected": false,
	}

	// Test the condition with the fact where it should be true
	satisfied, _, _, err := condition.Evaluate(factTrue, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected condition to be true, but it was false")
	}

	// Test the condition with the fact where it should be false
	satisfied, _, _, err = condition.Evaluate(factFalse, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if satisfied {
		t.Errorf("Expected condition to be false, but it was true")
	}
}

func TestConditionEvaluateInvalidOperator(t *testing.T) {
	// Define a condition with an invalid operator
	condition := Condition{
		Fact:     "temperature",
		Operator: "invalidOperator",
		Value:    30,
	}

	// Define a fact
	fact := Fact{
		"temperature": 40,
	}

	// Test the condition with the fact
	_, _, _, err := condition.Evaluate(fact, "Ignore")
	if err == nil {
		t.Errorf("Expected an error due to invalid operator, but got none")
	}
}

func TestConditionEvaluateUnexpectedValueType(t *testing.T) {
	// Define a condition where the Value field is of an unexpected type
	condition := Condition{
		Fact:     "temperature",
		Operator: "greaterThan",
		Value:    "thirty", // string instead of a number
	}

	// Define a fact
	fact := Fact{
		"temperature": 40,
	}

	// Test the condition with the fact
	_, _, _, err := condition.Evaluate(fact, "Ignore")
	if err == nil {
		t.Errorf("Expected an error due to unexpected value type, but got none")
	}
}

// TestConditionEvaluateContainsOperatorWithSlice tests the evaluation of a condition using the "contains" operator with a slice of strings.
//
// It defines a condition that checks if the "activities" fact contains the value "swimming". It then defines two facts, one where the condition should be true and another where it should be false. It evaluates the condition with both facts and verifies the expected results.
//
// Parameters:
// - t: the testing.T instance used for testing.
//
// Returns:
// - None.
func TestConditionEvaluateContainsOperatorWithSlice(t *testing.T) {
	// Define a condition that uses the "contains" operator with a slice of strings
	condition := Condition{
		Fact:     "activities",
		Operator: "contains",
		Value:    "swimming",
	}

	// Define a fact where the condition should be true
	factTrue := Fact{
		"activities": []string{"running", "swimming", "cycling"},
	}

	// Define a fact where the condition should be false
	factFalse := Fact{
		"activities": []string{"running", "cycling"},
	}

	// Test the condition with the fact where it should be true
	satisfied, _, _, err := condition.Evaluate(factTrue, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected condition to be true, but it was false")
	}

	// Test the condition with the fact where it should be false
	satisfied, _, _, err = condition.Evaluate(factFalse, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if satisfied {
		t.Errorf("Expected condition to be false, but it was true")
	}
}

// TestConditionEvaluateBoundary tests the evaluation of a condition at the boundary.
//
// This function initializes a condition with a fact and tests the evaluation of the condition
// with a fact exactly at the boundary. It checks if the condition is satisfied and if there
// are any errors. If the condition is satisfied, it fails the test and reports an error.
// The function does not return any values.
func TestConditionEvaluateBoundary(t *testing.T) {
	// Define a condition
	condition := Condition{
		Fact:     "temperature",
		Operator: "greaterThan",
		Value:    30,
	}

	// Define a fact exactly at the boundary
	factBoundary := Fact{
		"temperature": 30,
	}

	// Test the condition with the boundary fact
	satisfied, _, _, err := condition.Evaluate(factBoundary, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating condition: %v", err)
	}
	if satisfied {
		t.Errorf("Expected condition to be false at boundary, but it was true")
	}
}

func TestConditionEvaluateInvalidFactType(t *testing.T) {
	condition := Condition{
		Fact:     "temperature",
		Operator: "greaterThan",
		Value:    30,
	}

	// Passing a string instead of an integer
	factInvalid := Fact{
		"temperature": "hot",
	}

	_, _, _, err := condition.Evaluate(factInvalid, "Ignore")
	if err == nil {
		t.Errorf("Expected an error due to invalid fact type, but got none")
	}
}

func TestConditionEvaluateMissingFact(t *testing.T) {
	condition := Condition{
		Fact:     "humidity",
		Operator: "lessThan",
		Value:    50,
	}

	// Fact "humidity" is missing
	factMissing := Fact{
		"temperature": 35,
	}

	_, _, _, err := condition.Evaluate(factMissing, "Error")
	if err == nil {
		t.Errorf("Expected an error due to missing fact, but got none")
	}
}

func TestRuleEvaluateComplexNested(t *testing.T) {
	rule := Rule{
		Name:     "TestRule",
		Priority: 1,
		Conditions: Conditions{
			All: []Condition{
				{
					Any: []Condition{
						{
							Fact:     "temperature",
							Operator: "greaterThan",
							Value:    30,
						},
						{
							Fact:     "humidity",
							Operator: "lessThan",
							Value:    50,
						},
					},
				},
				{
					All: []Condition{
						{
							Fact:     "windSpeed",
							Operator: "equalTo",
							Value:    10,
						},
					},
				},
			},
		}}

	fact := Fact{
		"temperature": 35,
		"humidity":    45,
		"windSpeed":   10,
	}

	satisfied, err := rule.Evaluate(fact, true, "Ignore")
	if err != nil {
		t.Fatalf("Error evaluating rule: %v", err)
	}
	if !satisfied {
		t.Errorf("Expected rule to be true, but it was false")
	}
}
