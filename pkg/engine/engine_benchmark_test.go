package engine_test

import (
	"fmt"
	"testing"

	"github.com/rgehrsitz/rulegopher/pkg/engine"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

func BenchmarkEvaluate(b *testing.B) {
	// Create a new engine
	e := engine.NewEngine()

	// Add a rule to the engine
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
			EventType: "alert",
		},
	}
	err := e.AddRule(rule)
	if err != nil {
		b.Fatalf("Failed to add rule: %v", err)
	}

	// Create a large number of facts
	facts := make([]rules.Fact, b.N)
	for i := 0; i < b.N; i++ {
		facts[i] = rules.Fact{
			"temperature": i,
		}
	}

	// Measure the time it takes to evaluate all the facts
	b.ResetTimer()
	for _, fact := range facts {
		_, err := e.Evaluate(fact)
		if err != nil {
			b.Fatalf("Failed to evaluate fact: %v", err)
		}
	}
}

func BenchmarkEngine_EvaluateRules_Performance(b *testing.B) {
	engine := engine.NewEngine()

	// Create a large number of rules
	for i := 0; i < 10000; i++ {
		rule := rules.Rule{
			Name:     fmt.Sprintf("Rule %d", i),
			Priority: i,
			Conditions: rules.Conditions{
				All: []rules.Condition{
					{
						Fact:     "temperature",
						Operator: "greaterThan",
						Value:    i,
					},
				},
			},
			Event: rules.Event{
				EventType: "High Temperature",
			},
		}
		err := engine.AddRule(rule)
		if err != nil {
			b.Fatalf("Failed to add rule: %v", err)
		}
	}

	// Create a fact
	fact := rules.Fact{
		"temperature": 5000,
	}

	// Reset the timer to exclude the setup time
	b.ResetTimer()

	// Run the Evaluate function b.N times
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(fact)
		if err != nil {
			b.Fatalf("Failed to evaluate facts: %v", err)
		}
	}
}

func BenchmarkEngine_EvaluateRule_Performance(b *testing.B) {
	engine := engine.NewEngine()

	// Create a large number of conditions
	conditions := make([]rules.Condition, 10000)
	for i := 0; i < 10000; i++ {
		conditions[i] = rules.Condition{
			Fact:     "temperature",
			Operator: "greaterThan",
			Value:    i,
		}
	}

	// Create a rule with the large number of conditions
	rule := rules.Rule{
		Name:     "Large Rule",
		Priority: 1,
		Conditions: rules.Conditions{
			All: conditions,
		},
		Event: rules.Event{
			EventType: "High Temperature",
		},
	}

	err := engine.AddRule(rule)
	if err != nil {
		b.Fatalf("Failed to add rule: %v", err)
	}

	// Create a fact
	fact := rules.Fact{
		"temperature": 5000,
	}

	// Reset the timer to exclude the setup time
	b.ResetTimer()

	// Run the Evaluate function b.N times
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(fact)
		if err != nil {
			b.Fatalf("Failed to evaluate facts: %v", err)
		}
	}
}
