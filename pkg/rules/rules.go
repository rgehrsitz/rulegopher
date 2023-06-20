package rules

import (
	"reflect"
	"strconv"
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
	Fact     string      `json:"fact,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	All      []Condition `json:"all,omitempty"`
	Any      []Condition `json:"any,omitempty"`
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
	// If this is a simple condition, evaluate it based on the fact, operator, and value
	if c.Fact != "" && c.Operator != "" {
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
			// Convert the fact value and condition value to float64
			factFloat, ok1 := convertToFloat64(factValue)
			valueFloat, ok2 := convertToFloat64(c.Value)
			if !ok1 || !ok2 {
				return false
			}
			switch c.Operator {
			case "greaterThan":
				return factFloat > valueFloat
			case "greaterThanOrEqual":
				return factFloat >= valueFloat
			case "lessThan":
				return factFloat < valueFloat
			case "lessThanOrEqual":
				return factFloat <= valueFloat
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
	}
	// If this is a complex condition, evaluate the nested conditions
	for _, condition := range c.All {
		if !condition.Evaluate(fact) {
			// If any 'all' condition is not satisfied, the condition is not satisfied
			return false
		}
	}
	if len(c.Any) > 0 {
		for _, condition := range c.Any {
			if condition.Evaluate(fact) {
				// If any 'any' condition is satisfied, the condition is satisfied
				return true
			}
		}
		// If no 'any' conditions are satisfied, the condition is not satisfied
		return false
	}

	// If there are no 'all' or 'any' conditions, the condition is satisfied
	return true
}

func convertToFloat64(value interface{}) (float64, bool) {
	switch value := value.(type) {
	case int:
		return float64(value), true
	case float64:
		return value, true
	case string:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v, true
		}
	}
	return 0, false
}
