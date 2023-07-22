package rules

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Rule represents a rule with a name, priority, conditions, and an event.
// @property {string} Name - The Name property is a string that represents the name of the rule.
// @property {int} Priority - Priority is an integer value that determines the order in which rules are
// evaluated. Rules with higher priority values are evaluated before rules with lower priority values.
// @property {Conditions} Conditions - Conditions is a struct that represents the conditions that need
// to be met for the rule to be triggered. It contains properties such as "Field" (the field to be
// checked), "Operator" (the comparison operator), and "Value" (the value to compare against).
// @property {Event} Event - The `Event` property represents an event that triggers the rule. It could
// be any kind of event, such as a user action, a system event, or a time-based event.
type Rule struct {
	Name       string     `json:"name"`
	Priority   int        `json:"priority"`
	Conditions Conditions `json:"conditions"`
	Event      Event      `json:"event"`
}

// Event defines a struct type named "Event" with various fields and JSON tags.
// @property {string} EventType - The EventType property is a string that represents the type of event.
// It can be used to categorize different types of events.
// @property CustomProperty - The `CustomProperty` field is of type `interface{}`. This means it can
// hold values of any type. It is used to store custom properties or additional information related to
// the event.
// @property {[]string} Facts - Facts is a slice of strings that represents additional information or
// data related to the event. It is optional and can be omitted if not needed.
// @property {[]interface{}} Values - The `Values` property is an array of interface{} type, which
// means it can hold values of any type. It is marked as `omitempty`, which means it will be omitted
// from the JSON output if it is empty.
// @property {string} RuleName - The `RuleName` property is a string that represents the name of the
// rule associated with the event. It is an optional property and may be omitted if not applicable.
type Event struct {
	EventType      string        `json:"eventType"`
	CustomProperty interface{}   `json:"customProperty"`
	Facts          []string      `json:"facts,omitempty"`
	Values         []interface{} `json:"values,omitempty"`
	RuleName       string        `json:"ruleName,omitempty"`
}

// Conditions is a struct that contains two arrays of Condition structs, one for all
// conditions and one for any conditions.
// @property {[]Condition} All - The "All" property is an array of conditions. It represents a set of
// conditions that must all be true for a certain condition to be considered true. In other words, all
// conditions in the "All" array must evaluate to true for the overall condition to be true.
// @property {[]Condition} Any - The `Any` property is an array of `Condition` objects. It represents a
// list of conditions where at least one of them must be true for the overall condition to be true.
type Conditions struct {
	All []Condition `json:"all"`
	Any []Condition `json:"any"`
}

// Condition represents a condition with a fact, operator, value, and optional nested conditions.
// @property {string} Fact - The "Fact" property represents the specific fact or attribute that the
// condition is checking. It is a string that describes the fact being evaluated.
// @property {string} Operator - The "Operator" property in the "Condition" struct represents the
// comparison operator to be used in the condition. It specifies how the "Fact" property should be
// compared to the "Value" property. Examples of operators include "=", ">", "<", ">=", "<=", "!=",
// "in", "not
// @property Value - The `Value` property is used to store the value that will be compared with the
// fact using the specified operator. The type of the value can be any valid Go data type, such as
// string, number, boolean, or even a custom struct.
// @property {[]Condition} All - The `All` property is an array of `Condition` objects. It represents a
// logical AND condition, where all the conditions in the array must be true for the overall condition
// to be true.
// @property {[]Condition} Any - The `Any` property is an array of `Condition` objects. It represents a
// logical OR condition, where at least one of the conditions in the array must be true for the overall
// condition to be true.
type Condition struct {
	Fact     string      `json:"fact,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	All      []Condition `json:"all,omitempty"`
	Any      []Condition `json:"any,omitempty"`
}

// Fact is a map with string keys and interface{} values.
type Fact map[string]interface{}

// The line `const epsilon = 1e-9` is declaring a constant named `epsilon` with a value of `1e-9`. This
// constant is used as a small value to determine if two floating-point numbers are almost equal. It is
// used in the `almostEqual` function to check if the absolute difference between two numbers is less
// than or equal to `epsilon`.
const epsilon = 1e-9

// almostEqual checks if two floating-point numbers are almost equal within a certain
// epsilon value.
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
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
	satisfied, facts, values, err := evaluateConditions(r.Conditions.All, fact)
	if err != nil {
		return false, err
	}
	if !satisfied {
		return false, nil
	}
	if satisfied && includeTriggeringFact {
		event := r.Event
		event.Facts = append(event.Facts, facts...)
		event.Values = append(event.Values, values...)
		r.Event = event
	}

	if len(r.Conditions.Any) > 0 {
		satisfied, facts, values, err = evaluateConditions(r.Conditions.Any, fact)
		if err != nil {
			return false, err
		}
		if satisfied {
			if includeTriggeringFact {
				event := r.Event
				event.Facts = append(event.Facts, facts...)
				event.Values = append(event.Values, values...)
				r.Event = event
			}
			return true, nil
		}
	}

	return satisfied, nil
}

// Evaluate is a method of the `Condition` struct. It takes a `fact` of type `Fact` as a
// parameter and evaluates the condition against the given fact.
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

	if len(condition.All) > 0 {
		satisfied, facts, values, err := evaluateConditions(condition.All, fact)
		if err != nil {
			return false, nil, nil, err
		}
		if !satisfied {
			return false, nil, nil, nil
		}
		return true, facts, values, nil
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
		if !satisfied {
			return false, nil, nil, nil
		}
		return true, facts, values, nil
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
		satisfied, returnedFact, value, err := condition.Evaluate(fact)
		if err != nil {
			return false, nil, nil, err
		}
		if satisfied {
			facts = append(facts, returnedFact...)
			values = append(values, value...)
		}
	}

	if len(facts) == 0 {
		return false, nil, nil, nil
	}

	return true, facts, values, nil
}
