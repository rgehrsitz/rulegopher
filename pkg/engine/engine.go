package engine

import (
	"sort"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

// Engine represents a rule engine.
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

// AddRule adds a new rule to the rule engine.
func (e *Engine) AddRule(rule rules.Rule) error {
	if rule.Name == "" {
		return &EmptyRuleNameError{}
	}

	if err := rule.Validate(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if rule.Conditions.All == nil && rule.Conditions.Any == nil {
		return &NilRuleConditionsError{RuleName: rule.Name}
	}

	for _, existingRule := range e.Rules {
		if existingRule.Name == rule.Name {
			return &RuleAlreadyExistsError{RuleName: rule.Name}
		}
	}

	e.Rules[rule.Name] = rule
	e.addToIndex(&rule)

	return nil
}

// addToIndex adds a rule to the rule index.
func (e *Engine) addToIndex(rule *rules.Rule) {
	e.processConditions(rule.Conditions.All, rule)
	e.processConditions(rule.Conditions.Any, rule)
}

// processConditions processes conditions for indexing.
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

// insertRuleIntoIndex inserts a rule into the rule index.
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

// RemoveRule removes a rule from the rule engine.
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

// removeFromIndex removes a rule from the rule index.
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

// Evaluate evaluates the input fact against the rules.
func (e *Engine) Evaluate(inputFact rules.Fact) ([]rules.Event, error) {
	generatedEvents := make([]rules.Event, 0)
	evaluatedRules := make(map[string]bool) // Keep track of evaluated rules

	var matchingRules []*rules.Rule

	for factName := range inputFact {
		e.mu.RLock()
		if rules, ok := e.RuleIndex[factName]; ok {
			matchingRules = append(matchingRules, rules...)
		}
		e.mu.RUnlock()
	}

	var result *multierror.Error
	for _, rule := range matchingRules {
		if _, alreadyEvaluated := evaluatedRules[rule.Name]; !alreadyEvaluated {
			// Create a copy of the rule before evaluating it
			ruleCopy := *rule
			satisfied, err := ruleCopy.Evaluate(inputFact, e.ReportFacts)
			if err != nil {
				result = multierror.Append(result, err)
				continue
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

	return generatedEvents, result.ErrorOrNil()
}

// UpdateRule updates an existing rule in the rule engine.
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
