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
	Fact           string      `json:"fact,omitempty"`
	Value          interface{} `json:"value,omitempty"`
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
func (r *Rule) Evaluate(fact Fact, includeTriggeringFact bool) bool {
	// Evaluate the 'all' conditions
	for _, condition := range r.Conditions.All {
		satisfied, fact, value := condition.Evaluate(fact)
		if !satisfied {
			return false
		}
		if satisfied && includeTriggeringFact {
			r.Event.Fact = fact
			r.Event.Value = value
		}
	}

	// If there are no 'any' conditions, and all 'all' conditions are satisfied, the rule is satisfied
	if len(r.Conditions.Any) == 0 {
		return true
	}

	// Evaluate the 'any' conditions
	for _, condition := range r.Conditions.Any {
		satisfied, fact, value := condition.Evaluate(fact)
		if satisfied {
			if includeTriggeringFact {
				r.Event.Fact = fact
				r.Event.Value = value
			}
			return true
		}
	}

	// If no 'any' conditions are satisfied, the rule is not satisfied
	return false
}

func (c *Condition) Evaluate(fact Fact) (bool, string, interface{}) {
	// If this is a simple condition, evaluate it based on the fact, operator, and value
	if c.Fact != "" && c.Operator != "" {
		// Get the fact value
		factValue, ok := fact[c.Fact]
		if !ok {
			// If the fact is not present, the condition is not satisfied
			return false, "", nil
		}

		// Compare the fact value to the condition value based on the operator
		switch c.Operator {
		case "equal":
			if reflect.DeepEqual(factValue, c.Value) {
				return true, c.Fact, c.Value
			}
		case "notEqual":
			if !reflect.DeepEqual(factValue, c.Value) {
				return true, c.Fact, c.Value
			}
		case "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual":
			// Convert the fact value and condition value to float64
			factFloat, ok1 := convertToFloat64(factValue)
			valueFloat, ok2 := convertToFloat64(c.Value)
			if !ok1 || !ok2 {
				return false, "", nil
			}
			switch c.Operator {
			case "greaterThan":
				// assuming factValue and c.Value are float64 for simplicity
				if factFloat > valueFloat {
					return true, c.Fact, c.Value
				}
			case "greaterThanOrEqual":
				if factFloat >= valueFloat {
					return true, c.Fact, c.Value
				}
			case "lessThan":
				if factFloat < valueFloat {
					return true, c.Fact, c.Value
				}
			case "lessThanOrEqual":
				if factFloat <= valueFloat {
					return true, c.Fact, c.Value
				}
			}
		case "contains":
			// This operator is only supported for strings
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := c.Value.(string)
			if ok1 && ok2 && strings.Contains(factStr, valueStr) {
				return true, c.Fact, c.Value
			}
		case "notContains":
			// This operator is only supported for strings
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := c.Value.(string)
			if ok1 && ok2 && !strings.Contains(factStr, valueStr) {
				return true, c.Fact, c.Value
			}
		}

		// If the operator is not recognized or the condition is not satisfied, return false
		return false, "", nil
	}

	// If this is a complex condition, evaluate the nested conditions
	for _, condition := range c.All {
		satisfied, fact, value := condition.Evaluate(fact)
		if !satisfied {
			// If any 'all' condition is not satisfied, the condition is not satisfied
			return false, "", nil
		}
		if satisfied {
			return true, fact, value
		}
	}
	if len(c.Any) > 0 {
		for _, condition := range c.Any {
			satisfied, fact, value := condition.Evaluate(fact)
			if satisfied {
				return true, fact, value
			}
		}
		// If no 'any' conditions are satisfied, the condition is not satisfied
		return false, "", nil
	}

	// If there are no 'all' or 'any' conditions, the condition is satisfied
	return true, "", nil
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
