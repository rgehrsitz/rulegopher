package engine

import (
	"errors"
	"sort"
	"sync"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type Engine struct {
	Rules          []rules.Rule
	RuleIndex      map[string][]*rules.Rule
	mu             sync.RWMutex
	ReportFacts    bool
	ReportRuleName bool
}

func NewEngine() *Engine {
	return &Engine{
		Rules:          make([]rules.Rule, 0),
		RuleIndex:      make(map[string][]*rules.Rule),
		ReportFacts:    false,
		ReportRuleName: false,
	}
}

func (e *Engine) AddRule(rule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, existingRule := range e.Rules {
		if existingRule.Name == rule.Name {
			return errors.New("rule already exists")
		}
	}

	// Insert the rule in the correct position to maintain sorted order
	insertionIndex := sort.Search(len(e.Rules), func(i int) bool {
		return e.Rules[i].Priority > rule.Priority
	})
	e.Rules = append(e.Rules, rules.Rule{})
	copy(e.Rules[insertionIndex+1:], e.Rules[insertionIndex:])
	e.Rules[insertionIndex] = rule

	e.addToIndex(&rule)
	return nil
}

func (e *Engine) addToIndex(rule *rules.Rule) {
	for _, condition := range rule.Conditions.All {
		e.insertRuleIntoIndex(condition.Fact, rule)
	}
	for _, condition := range rule.Conditions.Any {
		e.insertRuleIntoIndex(condition.Fact, rule)
	}
}

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

func (e *Engine) RemoveRule(ruleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for ruleIndex, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules = append(e.Rules[:ruleIndex], e.Rules[ruleIndex+1:]...)
			e.removeFromIndex(&rule)
			return nil
		}
	}

	return errors.New("rule does not exist")
}

func (e *Engine) removeFromIndex(rule *rules.Rule) {
	for factName, matchingRules := range e.RuleIndex {
		for ruleIndex, r := range matchingRules {
			if r == rule {
				e.RuleIndex[factName] = append(matchingRules[:ruleIndex], matchingRules[ruleIndex+1:]...)
				break
			}
		}
	}
}

func (e *Engine) Evaluate(inputFact rules.Fact) []rules.Event {
	e.mu.RLock()
	defer e.mu.RUnlock()

	generatedEvents := make([]rules.Event, 0)
	evaluatedRules := make(map[*rules.Rule]bool)
	for factName := range inputFact {
		matchingRules, ok := e.RuleIndex[factName]
		if ok {
			for _, rule := range matchingRules {
				if _, alreadyEvaluated := evaluatedRules[rule]; !alreadyEvaluated {
					// Create a copy of the rule before evaluating it
					ruleCopy := *rule
					if ruleCopy.Evaluate(inputFact, e.ReportFacts) {
						if e.ReportRuleName { // Check if the ReportRuleName option is enabled
							ruleCopy.Event.RuleName = ruleCopy.Name // Set the RuleName field here
						}
						generatedEvents = append(generatedEvents, ruleCopy.Event)
					}
					evaluatedRules[rule] = true
				}
			}
		}
	}

	return generatedEvents
}

func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for ruleIndex, existingRule := range e.Rules {
		if existingRule.Name == ruleName {
			e.removeFromIndex(&existingRule)
			oldPriority := existingRule.Priority
			e.Rules[ruleIndex] = newRule
			e.addToIndex(&newRule)
			// Re-sort the rules after updating only if the priority has changed
			if oldPriority != newRule.Priority {
				sort.Slice(e.Rules, func(i, j int) bool {
					return e.Rules[i].Priority < e.Rules[j].Priority
				})
			}
			return nil
		}
	}

	return errors.New("rule does not exist")
}
