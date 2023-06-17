package engine

import (
	"errors"
	"sync"

	"github.com/rgehrsitz/rulegopher/pkg/rules"
)

type Engine struct {
	Rules []rules.Rule
	mu    sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		Rules: make([]rules.Rule, 0),
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
	return nil
}

func (e *Engine) RemoveRule(ruleName string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules = append(e.Rules[:i], e.Rules[i+1:]...)
			return nil
		}
	}

	return errors.New("rule does not exist")
}

func (e *Engine) Evaluate(fact rules.Fact) []rules.Event {
	e.mu.RLock()
	defer e.mu.RUnlock()

	events := make([]rules.Event, 0)
	for _, rule := range e.Rules {
		if rule.Evaluate(fact) {
			events = append(events, rule.Event)
		}
	}

	return events
}

func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules[i] = newRule
			return nil
		}
	}

	return errors.New("rule does not exist")
}
