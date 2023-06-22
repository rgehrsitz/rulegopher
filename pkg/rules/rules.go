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
	EventType      string        `json:"eventType"`
	CustomProperty interface{}   `json:"customProperty"`
	Facts          []string      `json:"facts,omitempty"`
	Values         []interface{} `json:"values,omitempty"`
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

func (r *Rule) Evaluate(fact Fact, includeTriggeringFact bool) bool {
	for _, condition := range r.Conditions.All {
		satisfied, facts, values := condition.Evaluate(fact)
		if !satisfied {
			return false
		}
		if satisfied && includeTriggeringFact {
			r.Event.Facts = append(r.Event.Facts, facts...)
			r.Event.Values = append(r.Event.Values, values...)
		}
	}

	if len(r.Conditions.Any) == 0 {
		return true
	}

	for _, condition := range r.Conditions.Any {
		satisfied, facts, values := condition.Evaluate(fact)
		if satisfied {
			if includeTriggeringFact {
				r.Event.Facts = append(r.Event.Facts, facts...)
				r.Event.Values = append(r.Event.Values, values...)
			}
			return true
		}
	}

	return false
}

func (c *Condition) Evaluate(fact Fact) (bool, []string, []interface{}) {
	if c.Fact != "" && c.Operator != "" {
		factValue, ok := fact[c.Fact]
		if !ok {
			return false, nil, nil
		}

		switch c.Operator {
		case "equal":
			if reflect.DeepEqual(factValue, c.Value) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}
		case "notEqual":
			if !reflect.DeepEqual(factValue, c.Value) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}
		case "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual":
			factFloat, ok1 := convertToFloat64(factValue)
			valueFloat, ok2 := convertToFloat64(c.Value)
			if !ok1 || !ok2 {
				return false, nil, nil
			}
			switch c.Operator {
			case "greaterThan":
				if factFloat > valueFloat {
					return true, []string{c.Fact}, []interface{}{factValue}
				}
			case "greaterThanOrEqual":
				if factFloat >= valueFloat {
					return true, []string{c.Fact}, []interface{}{factValue}
				}
			case "lessThan":
				if factFloat < valueFloat {
					return true, []string{c.Fact}, []interface{}{factValue}
				}
			case "lessThanOrEqual":
				if factFloat <= valueFloat {
					return true, []string{c.Fact}, []interface{}{factValue}
				}
			}
		case "contains":
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := c.Value.(string)
			if ok1 && ok2 && strings.Contains(factStr, valueStr) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}
			factSlice, ok3 := factValue.([]string)
			if ok3 && contains(factSlice, valueStr) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}
		case "notContains":
			factStr, ok1 := factValue.(string)
			valueStr, ok2 := c.Value.(string)
			if ok1 && ok2 && !strings.Contains(factStr, valueStr) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}
			factSlice, ok3 := factValue.([]string)
			if ok3 && !contains(factSlice, valueStr) {
				return true, []string{c.Fact}, []interface{}{factValue}
			}

		}
		return false, nil, nil
	}

	var facts []string
	var values []interface{}

	for _, condition := range c.All {
		satisfied, fact, value := condition.Evaluate(fact)
		if !satisfied {
			return false, nil, nil
		}
		if satisfied {
			facts = append(facts, fact...)
			values = append(values, value...)
		}
	}

	if len(c.Any) > 0 {
		for _, condition := range c.Any {
			satisfied, fact, value := condition.Evaluate(fact)
			if satisfied {
				facts = append(facts, fact...)
				values = append(values, value...)
			}
		}
		if len(facts) == 0 {
			return false, nil, nil
		}
	}

	return true, facts, values
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

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
