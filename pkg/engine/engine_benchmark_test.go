package engine_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

func BenchmarkLargeRuleset(b *testing.B) {

	numIterations := 100

	// Read rules.json into a string
	rulesJSON := readFile("rules.json")

	// Unmarshal into a Rules struct
	var ruleSet []rules.Rule
	json.Unmarshal([]byte(rulesJSON), &ruleSet)

	engine := engine.NewEngine()

	for _, rule := range ruleSet {
		engine.AddRule(rule)
	}

	// Read facts.json into a string
	factsJSON := readFile("facts.json")

	// Unmarshal into a []Fact
	var facts []rules.Fact
	json.Unmarshal([]byte(factsJSON), &facts)

	// Benchmark Evaluate
	start := time.Now()

	for i := 0; i < numIterations; i++ {
		for _, fact := range facts {
			engine.Evaluate(fact)
		}
	}

	// Benchmarking stats
	totalDuration := time.Since(start)

	fmt.Println("Total duration:", totalDuration)

	numRules := len(ruleSet)
	numFacts := len(facts)

	perIteration := totalDuration / time.Duration(numIterations)

	fmt.Println("Duration per iteration:", perIteration)

	// Print stats
	fmt.Println("Num rules:", numRules)
	fmt.Println("Num facts:", numFacts)
	fmt.Println("Num iterations:", numIterations)

	// Calculate facts/sec
	factsPerSec := float64(numFacts*numIterations) / totalDuration.Seconds()
	fmt.Println("Facts per second:", factsPerSec)

	// Calculate rules/sec
	rulesPerSec := float64(numRules*numIterations) / totalDuration.Seconds()
	fmt.Println("Rules per second:", rulesPerSec)

}

func readFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
