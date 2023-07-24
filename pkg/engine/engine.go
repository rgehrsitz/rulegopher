package engine

import (
	"sort"
	"sync"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

// Engine represents a rule engine with rules and related properties.
// @property Rules - The `Rules` property is a map that stores the rules of the engine. Each rule is
// identified by a string key and the corresponding value is an instance of the `rules.Rule` struct.
// @property RuleIndex - The RuleIndex property is a map that stores the rules indexed by their names.
// Each rule name is mapped to a slice of pointers to Rule objects. This allows for efficient lookup of
// rules by their names.
// @property mu - The `mu` property is a `sync.RWMutex` type, which is a mutual exclusion lock. It
// provides a way to synchronize access to shared resources by allowing multiple readers or a single
// writer at a time. In this case, it is used to protect concurrent access to the `Engine`
// @property {bool} ReportFacts - A boolean value indicating whether to report the facts during the
// execution of the rules.
// @property {bool} ReportRuleName - The `ReportRuleName` property is a boolean flag that determines
// whether or not to include the rule name in the generated report. If set to `true`, the rule name
// will be included in the report. If set to `false`, the rule name will be excluded from the report.
type Engine struct {
	Rules          map[string]rules.Rule
	RuleIndex      map[string][]*rules.Rule
	mu             sync.RWMutex
	ReportFacts    bool
	ReportRuleName bool
}

// NewEngine returns a new instance of the Engine struct with initialized maps.
func NewEngine() *Engine {
	return &Engine{
		Rules:          make(map[string]rules.Rule),
		RuleIndex:      make(map[string][]*rules.Rule),
		ReportFacts:    false,
		ReportRuleName: false,
	}
}

// AddRule is a method of the `Engine` struct. It adds a new rule to the rule engine.
func (e *Engine) AddRule(rule rules.Rule) error {
	// Check if the rule name is empty
	if rule.Name == "" {
		return &EmptyRuleNameError{}
	}

	// Validate the rule before adding it
	if err := rule.Validate(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the Conditions field of the rule is nil
	if rule.Conditions.All == nil && rule.Conditions.Any == nil {
		return &NilRuleConditionsError{RuleName: rule.Name}
	}

	// Check if the rule already exists
	for _, existingRule := range e.Rules {
		if existingRule.Name == rule.Name {
			return &RuleAlreadyExistsError{RuleName: rule.Name}
		}
	}

	e.Rules[rule.Name] = rule
	e.addToIndex(&rule)

	return nil
}

// addToIndex is a method of the `Engine` struct and is responsible for adding a rule to the rule index.
// It iterates over the conditions of the rule and calls the `insertRuleIntoIndex` method for each
// condition. This ensures that the rule is correctly indexed based on its conditions.
func (e *Engine) addToIndex(rule *rules.Rule) {
	e.processConditions(rule.Conditions.All, rule)
	e.processConditions(rule.Conditions.Any, rule)
}

func (e *Engine) processConditions(conditions []rules.Condition, rule *rules.Rule) {
	for _, condition := range conditions {
		e.insertRuleIntoIndex(condition.Fact, rule)
		if len(condition.All) > 0 {
			e.processConditions(condition.All, rule)
		}
		if len(condition.Any) > 0 {
			e.processConditions(condition.Any, rule)
		}
	}
}

// insertRuleIntoIndex is responsible for inserting a rule into the rule index of the
// engine.
func (e *Engine) insertRuleIntoIndex(fact string, rule *rules.Rule) {
	existingRules := e.RuleIndex[fact]

	// Find the correct position to insert the new rule
	insertionIndex := sort.Search(len(existingRules), func(i int) bool {
		return existingRules[i].Priority > rule.Priority
	})
	// Insert the rule in the correct position
	existingRules = append(existingRules, nil)
	copy(existingRules[insertionIndex+1:], existingRules[insertionIndex:])
	existingRules[insertionIndex] = rule
	e.RuleIndex[fact] = existingRules
}

// RemoveRule is a method of the `Engine` struct. It is used to remove a rule from the
// rule engine.
func (e *Engine) RemoveRule(ruleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the rule exists
	if _, exists := e.Rules[ruleName]; !exists {
		return &RuleDoesNotExistError{RuleName: ruleName}
	}

	delete(e.Rules, ruleName)
	e.removeFromIndex(ruleName)

	return nil
}

// removeFromIndex is a method of the `Engine` struct and is responsible for removing a rule from the rule
// index. It iterates over the rule index and checks if the rule name matches the given `ruleName`. If
// a match is found, it removes the rule from the slice of matching rules for that fact. It does this
// by using the `append` function to create a new slice that excludes the rule at the specified index.
func (e *Engine) removeFromIndex(ruleName string) {
	for factName, matchingRules := range e.RuleIndex {
		for ruleIndex, r := range matchingRules {
			if r.Name == ruleName {
				e.RuleIndex[factName] = append(matchingRules[:ruleIndex], matchingRules[ruleIndex+1:]...)
				break
			}
		}
	}
}

// Evaluate is a method of the `Engine` struct and is responsible for evaluating the input fact against
// the rules in the rule engine.
func (e *Engine) Evaluate(inputFact rules.Fact) ([]rules.Event, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	generatedEvents := make([]rules.Event, 0)
	evaluatedRules := make(map[string]bool) // Keep track of evaluated rules
	for factName := range inputFact {
		matchingRules, ok := e.RuleIndex[factName]
		if ok {
			for _, rule := range matchingRules {
				if _, alreadyEvaluated := evaluatedRules[rule.Name]; !alreadyEvaluated {
					// Create a copy of the rule before evaluating it
					ruleCopy := *rule
					satisfied, err := ruleCopy.Evaluate(inputFact, e.ReportFacts)
					if err != nil {
						return nil, err
					}
					if satisfied {
						if e.ReportRuleName { // Check if the ReportRuleName option is enabled
							ruleCopy.Event.RuleName = ruleCopy.Name // Set the RuleName field here
						}
						generatedEvents = append(generatedEvents, ruleCopy.Event)
					}
					evaluatedRules[rule.Name] = true
				}
			}
		}
	}

	return generatedEvents, nil
}

// UpdateRule is a function of the `Engine` struct. It is used to update an existing rule
// in the rule engine.
func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Validate the new rule before updating
	if err := newRule.Validate(); err != nil {
		return &InvalidRuleError{RuleName: newRule.Name}
	}

	// Check if the rule exists
	if _, exists := e.Rules[ruleName]; !exists {
		return &RuleDoesNotExistError{RuleName: ruleName}
	}

	e.removeFromIndex(ruleName)
	e.Rules[ruleName] = newRule
	e.addToIndex(&newRule)

	return nil
}
