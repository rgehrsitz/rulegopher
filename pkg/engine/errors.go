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
