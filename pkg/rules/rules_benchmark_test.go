package rules

import (
	"testing"
)

func BenchmarkEvaluate(b *testing.B) {
	// Initialize a rule and a fact
	rule := Rule{
		Name:     "Rule1",
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
			EventType:      "Event1",
			CustomProperty: "Custom1",
		},
	}
	fact := Fact{
		"temperature": 35,
	}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Evaluate(fact, true)
	}
}
