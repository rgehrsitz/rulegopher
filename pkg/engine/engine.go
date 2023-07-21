package engine

import (
	"errors"
	"sort"
	"sync"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type Engine struct {
	Rules          map[string]rules.Rule
	RuleIndex      map[string][]*rules.Rule
	mu             sync.RWMutex
	ReportFacts    bool
	ReportRuleName bool
}

func NewEngine() *Engine {
	return &Engine{
		Rules:          make(map[string]rules.Rule),
		RuleIndex:      make(map[string][]*rules.Rule),
		ReportFacts:    false,
		ReportRuleName: false,
	}
}

func (e *Engine) AddRule(rule rules.Rule) error {
	// Check if the rule name is empty
	if rule.Name == "" {
		return errors.New("rule name cannot be empty")
	}

	// Validate the rule before adding it
	if err := rule.Validate(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Check if the Conditions field of the rule is nil
	if rule.Conditions.All == nil && rule.Conditions.Any == nil {
		return errors.New("rule conditions cannot be nil")
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

func (e *Engine) addToIndex(rule *rules.Rule) {
	for _, condition := range rule.Conditions.All {
		e.insertRuleIntoIndex(condition.Fact, rule.Name)
	}
	for _, condition := range rule.Conditions.Any {
		e.insertRuleIntoIndex(condition.Fact, rule.Name)
	}
}

func (e *Engine) insertRuleIntoIndex(fact string, ruleName string) {
	existingRules := e.RuleIndex[fact]

	// Find the rule by its name
	var rule *rules.Rule
	for _, r := range e.Rules {
		if r.Name == ruleName {
			rule = &r
			break
		}
	}
	if rule == nil {
		return
	}

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

func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Validate the new rule before updating
	if err := newRule.Validate(); err != nil {
		return err
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
