package engine

import (
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

func (e *Engine) AddRule(rule rules.Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Rules = append(e.Rules, rule)
}

func (e *Engine) RemoveRule(ruleName string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules = append(e.Rules[:i], e.Rules[i+1:]...)
			break
		}
	}
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

func (e *Engine) UpdateRule(ruleName string, newRule rules.Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, rule := range e.Rules {
		if rule.Name == ruleName {
			e.Rules[i] = newRule
			break
		}
	}
}
