package engine_test

import (
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
