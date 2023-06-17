package rules

import (
	"reflect"
	"strings"
)

type Rule struct {
	Name       string     `json:"name"`
	Priority   int        `json:"priority"`
	Conditions Conditions `json:"conditions"`
	Event      Event      `json:"event"`
}

type Event struct {
	EventType      string      `json:"eventType"`
	CustomProperty interface{} `json:"customProperty"`
}

type Conditions struct {
	All []Condition `json:"all"`
	Any []Condition `json:"any"`
}

type Condition struct {
	Fact     string      `json:"fact"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Fact map[string]interface{}

// Add functions for adding, removing, and updating rules here.
func (r *Rule) Evaluate(fact Fact) bool {
	// Evaluate the 'all' conditions
	for _, condition := range r.Conditions.All {
		if !condition.Evaluate(fact) {
			// If any 'all' condition is not satisfied, the rule is not satisfied
			return false
		}
	}

	// If there are no 'any' conditions, and all 'all' conditions are satisfied, the rule is satisfied
	if len(r.Conditions.Any) == 0 {
		return true
	}

	// Evaluate the 'any' conditions
	for _, condition := range r.Conditions.Any {
		if condition.Evaluate(fact) {
			// If any 'any' condition is satisfied, the rule is satisfied
			return true
		}
	}

	// If no 'any' conditions are satisfied, the rule is not satisfied
	return false
}

func (c *Condition) Evaluate(fact Fact) bool {
	// Get the fact value
	factValue, ok := fact[c.Fact]
	if !ok {
		// If the fact is not present, the condition is not satisfied
		return false
	}

	// Compare the fact value to the condition value based on the operator
	switch c.Operator {
	case "equal":
		return reflect.DeepEqual(factValue, c.Value)
	case "notEqual":
		return !reflect.DeepEqual(factValue, c.Value)
	case "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual":
		// These operators are not supported for all types, so we'll handle each type separately
		switch factValue := factValue.(type) {
		case int:
			value, ok := c.Value.(int)
			if !ok {
				return false
			}
			switch c.Operator {
			case "greaterThan":
				return factValue > value
			case "greaterThanOrEqual":
				return factValue >= value
			case "lessThan":
				return factValue < value
			case "lessThanOrEqual":
				return factValue <= value
			}
		case float64:
			value, ok := c.Value.(float64)
			if !ok {
				return false
			}
			switch c.Operator {
			case "greaterThan":
				return factValue > value
			case "greaterThanOrEqual":
				return factValue >= value
			case "lessThan":
				return factValue < value
			case "lessThanOrEqual":
				return factValue <= value
			}
		default:
			// If the fact value is not a numeric type, these operators are not supported
			return false
		}
	case "contains":
		// This operator is only supported for strings
		factStr, ok1 := factValue.(string)
		valueStr, ok2 := c.Value.(string)
		if ok1 && ok2 {
			return strings.Contains(factStr, valueStr)
		}
		return false
	case "notContains":
		// This operator is only supported for strings
		factStr, ok1 := factValue.(string)
		valueStr, ok2 := c.Value.(string)
		if ok1 && ok2 {
			return !strings.Contains(factStr, valueStr)
		}
		return false
	default:
		// If the operator is not recognized, the condition is not satisfied
		return false
	}
	return false
}
