package rules

import (
	"testing"
)

func BenchmarkRuleEvaluateAnyConditions(b *testing.B) {
	// Setup: Create a sample rule with ANY conditions
	rule := Rule{
		Conditions: Conditions{
			Any: []Condition{
				{Fact: "age", Operator: "lessThan", Value: 18},
				{Fact: "status", Operator: "equal", Value: "inactive"},
			},
		},
	}

	fact := Fact{"age": 25, "status": "active"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rule.Evaluate(fact, true)
	}
}

func BenchmarkRuleEvaluateAllConditions(b *testing.B) {
	// Setup: Create a sample rule with ALL conditions
	rule := Rule{
		Conditions: Conditions{
			All: []Condition{
				{Fact: "age", Operator: "greaterThan", Value: 18},
				{Fact: "status", Operator: "equal", Value: "active"},
			},
		},
	}

	fact := Fact{"age": 25, "status": "active"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rule.Evaluate(fact, true)
	}
}

func BenchmarkRuleEvaluateMixedConditions(b *testing.B) {
	// Setup: Create a sample rule with mixed ALL and ANY conditions
	rule := Rule{
		Conditions: Conditions{
			All: []Condition{
				{Fact: "age", Operator: "greaterThan", Value: 18},
				{Fact: "status", Operator: "equal", Value: "active"},
			},
			Any: []Condition{
				{Fact: "country", Operator: "equal", Value: "USA"},
				{Fact: "membership", Operator: "contains", Value: "premium"},
			},
		},
	}

	fact := Fact{"age": 25, "status": "active", "country": "USA", "membership": "basic"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rule.Evaluate(fact, true)
	}
}

func BenchmarkEvaluateAllConditions(b *testing.B) {
	// Setup: Create a sample array of ALL conditions
	conditions := []Condition{
		{Fact: "age", Operator: "greaterThan", Value: 18},
		{Fact: "status", Operator: "equal", Value: "active"},
	}

	fact := Fact{"age": 25, "status": "active"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = evaluateConditions(conditions, fact)
	}
}

func BenchmarkEvaluateAnyConditions(b *testing.B) {
	// Setup: Create a sample array of ANY conditions
	conditions := []Condition{
		{Fact: "age", Operator: "lessThan", Value: 18},
		{Fact: "status", Operator: "equal", Value: "inactive"},
	}

	fact := Fact{"age": 25, "status": "active"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = evaluateConditions(conditions, fact)
	}
}

func BenchmarkEvaluateMixedConditions(b *testing.B) {
	// Setup: Create a sample array of mixed conditions (nested)
	conditions := []Condition{
		{Fact: "age", Operator: "greaterThan", Value: 18},
		{Fact: "status", Operator: "equal", Value: "active"},
		{Fact: "country", Operator: "equal", Value: "USA", All: []Condition{
			{Fact: "membership", Operator: "contains", Value: "premium"},
		}},
	}

	fact := Fact{"age": 25, "status": "active", "country": "USA", "membership": "basic"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = evaluateConditions(conditions, fact)
	}
}
