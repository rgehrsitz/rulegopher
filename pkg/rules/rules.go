package rules

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Rule represents a single rule with its conditions and event.
type Rule struct {
	Name       string     `json:"name"`
	Priority   int        `json:"priority"`
	Conditions Conditions `json:"conditions"`
	Event      Event      `json:"event"`
}

// Event represents the event that is generated when a rule is satisfied.
type Event struct {
	EventType      string        `json:"eventType"`
	CustomProperty interface{}   `json:"customProperty"`
	Facts          []string      `json:"facts,omitempty"`
	Values         []interface{} `json:"values,omitempty"`
	RuleName       string        `json:"ruleName,omitempty"`
}

// Conditions represents the conditions of a rule.
type Conditions struct {
	All []Condition `json:"all"`
	Any []Condition `json:"any"`
}

// Condition represents a single condition of a rule.
type Condition struct {
	Fact     string      `json:"fact,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	All      []Condition `json:"all,omitempty"`
	Any      []Condition `json:"any,omitempty"`
}

// Fact represents a fact that is evaluated against the conditions of a rule.
type Fact map[string]interface{}

const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}

func (r *Rule) Evaluate(fact Fact, includeTriggeringFact bool) (bool, error) {
	satisfied, facts, values, err := evaluateConditions(r.Conditions.All, fact)
	if err != nil {
		return false, err
	}
	if !satisfied {
		return false, nil
	}
	if satisfied && includeTriggeringFact {
		r.Event.Facts = append(r.Event.Facts, facts...)
		r.Event.Values = append(r.Event.Values, values...)
	}

	if len(r.Conditions.Any) > 0 {
		satisfied, facts, values, err = evaluateConditions(r.Conditions.Any, fact)
		if err != nil {
			return false, err
		}
		if satisfied {
			if includeTriggeringFact {
				r.Event.Facts = append(r.Event.Facts, facts...)
				r.Event.Values = append(r.Event.Values, values...)
			}
			return true, nil
		}
	}

	return len(r.Conditions.Any) == 0, nil
}

// Evaluate evaluates the condition against the given fact.
func (condition *Condition) Evaluate(fact Fact) (bool, []string, []interface{}, error) {
	validOperators := map[string]bool{
		"equal":              true,
		"notEqual":           true,
		"greaterThan":        true,
		"greaterThanOrEqual": true,
		"lessThan":           true,
		"lessThanOrEqual":    true,
		"contains":           true,
		"notContains":        true,
	}

	if _, ok := validOperators[condition.Operator]; !ok {
		return false, nil, nil, fmt.Errorf("invalid operator: %s", condition.Operator)
	}

	if condition.Fact != "" && condition.Operator != "" {
		factValue, ok := fact[condition.Fact]
		if !ok {
			return false, nil, nil, nil
		}

		switch condition.Operator {
		case "equal":
			if reflect.DeepEqual(factValue, condition.Value) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
		case "notEqual":
			if !reflect.DeepEqual(factValue, condition.Value) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
		case "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual":
			factFloat, _, err1 := convertToFloat64(factValue)
			valueFloat, _, err2 := convertToFloat64(condition.Value)
			if err1 != nil {
				return false, nil, nil, fmt.Errorf("error converting fact value to float64: %w", err1)
			}
			if err2 != nil {
				return false, nil, nil, fmt.Errorf("error converting condition value to float64: %w", err2)
			}
			switch condition.Operator {
			case "greaterThan":
				if factFloat > valueFloat+epsilon {
					return true, []string{condition.Fact}, []interface{}{factValue}, nil
				}
			case "greaterThanOrEqual":
				if almostEqual(factFloat, valueFloat) || factFloat > valueFloat {
					return true, []string{condition.Fact}, []interface{}{factValue}, nil
				}
			case "lessThan":
				if factFloat < valueFloat-epsilon {
					return true, []string{condition.Fact}, []interface{}{factValue}, nil
				}
			case "lessThanOrEqual":
				if almostEqual(factFloat, valueFloat) || factFloat < valueFloat {
					return true, []string{condition.Fact}, []interface{}{factValue}, nil
				}
			}
		case "contains":
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := condition.Value.(string)
			if ok1 && ok2 && strings.Contains(factStr, valueStr) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
			factSlice, ok3 := factValue.([]string)
			if ok3 && contains(factSlice, valueStr) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
		case "notContains":
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := condition.Value.(string)
			if ok1 && ok2 && !strings.Contains(factStr, valueStr) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
			factSlice, ok3 := factValue.([]string)
			if ok3 && !contains(factSlice, valueStr) {
				return true, []string{condition.Fact}, []interface{}{factValue}, nil
			}
		}
		return false, nil, nil, nil
	}

	return evaluateConditions(condition.All, fact)
}

// convertToFloat64 attempts to convert the given value to a float64.
func convertToFloat64(value interface{}) (float64, bool, error) {
	switch value := value.(type) {
	case int:
		return float64(value), true, nil
	case float64:
		return value, true, nil
	case string:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v, true, nil
		} else {
			return 0, false, err
		}
	}
	return 0, false, fmt.Errorf("unsupported type: %T", value)
}

// contains checks if the given slice contains the given string.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func evaluateConditions(conditions []Condition, fact Fact) (bool, []string, []interface{}, error) {
	var facts []string
	var values []interface{}

	for _, condition := range conditions {
		satisfied, fact, value, err := condition.Evaluate(fact)
		if err != nil {
			return false, nil, nil, err
		}
		if satisfied {
			facts = append(facts, fact...)
			values = append(values, value...)
		}
	}

	if len(facts) == 0 {
		return false, nil, nil, nil
	}

	return true, facts, values, nil
}
