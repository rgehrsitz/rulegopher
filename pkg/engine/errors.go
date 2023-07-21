package engine

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
