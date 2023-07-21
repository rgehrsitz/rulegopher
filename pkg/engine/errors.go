package engine

// The below code defines custom error types for different rule-related scenarios in Go.
// @property {string} RuleName - The RuleName property is a string that represents the name of a rule.
// It is used in the error messages to provide more information about the specific rule that caused the
// error.
type RuleAlreadyExistsError struct {
	RuleName string
}

func (e *RuleAlreadyExistsError) Error() string {
	return "Rule already exists: " + e.RuleName
}

type RuleDoesNotExistError struct {
	RuleName string
}

func (e *RuleDoesNotExistError) Error() string {
	return "Rule does not exist: " + e.RuleName
}

type InvalidRuleError struct {
	RuleName string
}

func (e *InvalidRuleError) Error() string {
	return "Invalid rule: " + e.RuleName
}

type EmptyRuleNameError struct{}

func (e *EmptyRuleNameError) Error() string {
	return "rule name cannot be empty"
}

type NilRuleConditionsError struct {
	RuleName string
}

func (e *NilRuleConditionsError) Error() string {
	return "rule conditions cannot be nil for rule: " + e.RuleName
}
