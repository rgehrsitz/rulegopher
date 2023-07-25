package rules

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Rule represents a rule with a name, priority, conditions, and an event.
type Rule struct {
	Name       string     `json:"name"`
	Priority   int        `json:"priority"`
	Conditions Conditions `json:"conditions"`
	Event      Event      `json:"event"`
}

// Event defines a struct type named "Event" with various fields and JSON tags.
type Event struct {
	EventType      string        `json:"eventType"`
	CustomProperty interface{}   `json:"customProperty"`
	Facts          []string      `json:"facts,omitempty"`
	Values         []interface{} `json:"values,omitempty"`
	RuleName       string        `json:"ruleName,omitempty"`
}

// Conditions is a struct that contains two arrays of Condition structs, one for all
// conditions and one for any conditions.
type Conditions struct {
	All []Condition `json:"all"`
	Any []Condition `json:"any"`
}

// Condition represents a condition with a fact, operator, value, and optional nested conditions.
type Condition struct {
	Fact     string      `json:"fact,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	All      []Condition `json:"all,omitempty"`
	Any      []Condition `json:"any,omitempty"`
}

// Fact is a map with string keys and interface{} values.
type Fact map[string]interface{}

// Constant used as a small value to determine if two floating-point numbers are almost equal. It is
// used in the `almostEqual` function to check if the absolute difference between two numbers is less
// than or equal to `epsilon`.
const epsilon = 1e-9

// almostEqual checks if two floating-point numbers are almost equal, considering both absolute and relative differences.
func almostEqual(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff <= epsilon {
		// handle the case of small numbers
		return true
	}
	// use relative error
	return diff <= epsilon*math.Max(math.Abs(a), math.Abs(b))
}

// Validate is a method of the `Rule` struct. It is used to validate the operators used
// in the conditions of the rule.
func (r *Rule) Validate() error {
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

	for _, condition := range r.Conditions.All {
		if len(condition.All) > 0 || len(condition.Any) > 0 {
			// This is a nested condition, so we don't need to validate the operator
			continue
		}
		if _, ok := validOperators[condition.Operator]; !ok {
			return fmt.Errorf("invalid operator: %s", condition.Operator)
		}
	}

	for _, condition := range r.Conditions.Any {
		if len(condition.All) > 0 || len(condition.Any) > 0 {
			// This is a nested condition, so we don't need to validate the operator
			continue
		}
		if _, ok := validOperators[condition.Operator]; !ok {
			return fmt.Errorf("invalid operator: %s", condition.Operator)
		}
	}

	return nil
}

// Evaluate is a method of the `Rule` struct. It takes a `fact` of type `Fact` and a
// boolean `includeTriggeringFact` as parameters.
func (r *Rule) Evaluate(fact Fact, includeTriggeringFact bool) (bool, error) {
	allSatisfied, facts, values, err := evaluateConditions(r.Conditions.All, fact)
	if err != nil {
		return false, err
	}
	if !allSatisfied && len(r.Conditions.All) > 0 {
		return false, nil
	}
	if allSatisfied && includeTriggeringFact {
		event := r.Event
		event.Facts = append(event.Facts, facts...)
		event.Values = append(event.Values, values...)
		r.Event = event
	}

	anySatisfied, facts, values, err := evaluateConditions(r.Conditions.Any, fact)
	if err != nil {
		return false, err
	}
	if anySatisfied {
		if includeTriggeringFact {
			event := r.Event
			event.Facts = append(event.Facts, facts...)
			event.Values = append(event.Values, values...)
			r.Event = event
		}
		return true, nil
	}

	return len(r.Conditions.Any) == 0 && allSatisfied, nil
}

// Evaluate is a method of the `Condition` struct. It takes a `fact` of type `Fact` as a
// parameter and evaluates the condition against the given fact.
func (condition *Condition) Evaluate(fact Fact) (bool, []string, []interface{}, error) {
	if len(condition.All) > 0 || len(condition.Any) > 0 {
		return condition.evaluateNestedConditions(fact)
	}

	return condition.evaluateSimpleCondition(fact)
}

// evaluateSimpleCondition evaluates a simple condition (i.e., a condition without nested conditions)
// and returns whether the condition is satisfied, along with the corresponding fact and value.
func (condition *Condition) evaluateSimpleCondition(fact Fact) (bool, []string, []interface{}, error) {
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

	if len(condition.All) > 0 {
		satisfied, facts, values, err := evaluateConditions(condition.All, fact)
		if err != nil {
			return false, nil, nil, err
		}
		if satisfied {
			return true, facts, values, nil
		}
	}

	if len(condition.Any) > 0 {
		satisfied, facts, values, err := evaluateConditions(condition.Any, fact)
		if err != nil {
			return false, nil, nil, err
		}
		if satisfied {
			return true, facts, values, nil
		}
	}

	return false, nil, nil, nil
}

// evaluateNestedConditions evaluates nested conditions and returns whether any conditions
// are satisfied, along with the corresponding facts and values.
func (condition *Condition) evaluateNestedConditions(fact Fact) (bool, []string, []interface{}, error) {
	satisfied, facts, values, err := evaluateConditions(condition.All, fact)
	if err != nil {
		return false, nil, nil, err
	}
	if satisfied {
		return true, facts, values, nil
	}

	satisfied, facts, values, err = evaluateConditions(condition.Any, fact)
	if err != nil {
		return false, nil, nil, err
	}
	if satisfied {
		return true, facts, values, nil
	}

	return false, nil, nil, nil
}

// convertToFloat64 takes in a value of any type and attempts to convert it to a
// float64, returning the converted value, a boolean indicating success or failure, and an error if
// applicable.
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

// contains checks if a given string is present in a slice of strings.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// evaluateConditions evaluates a list of conditions against a given fact and returns whether any conditions
// are satisfied, along with the corresponding facts and values.
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
			return true, facts, values, nil // return true as soon as a condition is satisfied
		}
	}

	return false, nil, nil, nil
}
