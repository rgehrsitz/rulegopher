package rules

import (
	"testing"
)

func BenchmarkEvaluate(b *testing.B) {
	// Setup
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
		b.Fatalf("Failed to add rule: %v", err)
	}

	fact := rules.Fact{
		"temperature": 35,
		"humidity":    0.4,
	}

	// Run the benchmark
	b.ResetTimer() // Reset the timer to ignore the setup time
	for i := 0; i < b.N; i++ {
		_, err := engine.Evaluate(fact)
		if err != nil {
			b.Fatalf("Failed to evaluate facts: %v", err)
		}
	}
}
