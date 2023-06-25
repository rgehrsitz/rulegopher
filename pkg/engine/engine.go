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
	index := sort.Search(len(e.Rules), func(i int) bool {
		return e.Rules[i].Priority > rule.Priority
	})
	e.Rules = append(e.Rules, rules.Rule{})
	copy(e.Rules[index+1:], e.Rules[index:])
	e.Rules[index] = rule

	e.addToIndex(&rule)
	return nil
}

func (e *Engine) addToIndex(rule *rules.Rule) {
	for _, condition := range rule.Conditions.All {
		rules := e.RuleIndex[condition.Fact]
		// Find the correct position to insert the new rule
		index := sort.Search(len(rules), func(i int) bool {
			return rules[i].Priority > rule.Priority
		})
		// Insert the rule in the correct position
		rules = append(rules, nil)
		copy(rules[index+1:], rules[index:])
		rules[index] = rule
		e.RuleIndex[condition.Fact] = rules
	}
	for _, condition := range rule.Conditions.Any {
		rules := e.RuleIndex[condition.Fact]
		// Find the correct position to insert the new rule
		index := sort.Search(len(rules), func(i int) bool {
			return rules[i].Priority > rule.Priority
		})
		// Insert the rule in the correct position
		rules = append(rules, nil)
		copy(rules[index+1:], rules[index:])
		rules[index] = rule
		e.RuleIndex[condition.Fact] = rules
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
					if ruleCopy.Evaluate(fact, e.ReportFacts) {
						if e.ReportRuleName { // Check if the ReportRuleName option is enabled
							ruleCopy.Event.RuleName = ruleCopy.Name // Set the RuleName field here
						}
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
			// Re-sort the rules after updating
			sort.Slice(e.Rules, func(i, j int) bool {
				return e.Rules[i].Priority < e.Rules[j].Priority
			})
			return nil
		}
	}

	return errors.New("rule does not exist")
}
