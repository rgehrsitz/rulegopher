package engine

import (
	"errors"
	"sync"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type Engine struct {
	Rules     []rules.Rule
	RuleIndex map[string][]*rules.Rule
	mu        sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		Rules:     make([]rules.Rule, 0),
		RuleIndex: make(map[string][]*rules.Rule),
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

	e.Rules = append(e.Rules, rule)
	e.addToIndex(&rule)
	return nil
}

func (e *Engine) addToIndex(rule *rules.Rule) {
	for _, condition := range rule.Conditions.All {
		e.RuleIndex[condition.Fact] = append(e.RuleIndex[condition.Fact], rule)
	}
	for _, condition := range rule.Conditions.Any {
		e.RuleIndex[condition.Fact] = append(e.RuleIndex[condition.Fact], rule)
	}
}

func (e *Engine) RemoveRule(ruleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules = append(e.Rules[:i], e.Rules[i+1:]...)
			e.removeFromIndex(&rule)
			return nil
		}
	}

	return errors.New("rule does not exist")
}

func (e *Engine) removeFromIndex(rule *rules.Rule) {
	for fact, rules := range e.RuleIndex {
		for i, r := range rules {
			if r == rule {
				e.RuleIndex[fact] = append(rules[:i], rules[i+1:]...)
				break
			}
		}
	}
}

func (e *Engine) Evaluate(fact rules.Fact) []rules.Event {
	e.mu.RLock()
	defer e.mu.RUnlock()

	events := make([]rules.Event, 0)
	evaluatedRules := make(map[*rules.Rule]bool)
	for factKey := range fact {
		rules, ok := e.RuleIndex[factKey]
		if ok {
			for _, rule := range rules {
				if _, evaluated := evaluatedRules[rule]; !evaluated {
					// Create a copy of the rule before evaluating it
					ruleCopy := *rule
					if ruleCopy.Evaluate(fact, true) {
						events = append(events, ruleCopy.Event)
					}
					evaluatedRules[rule] = true
				}
			}
		}
	}

	return events
}

func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.removeFromIndex(&rule)
			e.Rules[i] = newRule
			e.addToIndex(&newRule)
			return nil
		}
	}

	return errors.New("rule does not exist")
}
