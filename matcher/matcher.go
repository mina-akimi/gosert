package matcher

import (
	"fmt"
)

// SuccessMatcherInstance is a singleton of SuccessMatcher.
var SuccessMatcherInstance = &SuccessMatcher{}

// SuccessMatcher always succeeds.  It is used as a placeholder for success matches.
type SuccessMatcher struct {
}

// Match implements types.GomegaMatcher.
func (matcher *SuccessMatcher) Match(actual interface{}) (success bool, err error) {
	return true, nil
}

// FailureMessage implements types.GomegaMatcher.
func (matcher *SuccessMatcher) FailureMessage(actual interface{}) (message string) {
	return ""
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (matcher *SuccessMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return ""
}

// FailureMatcher always fails.  It is used to report custom error messages.
type FailureMatcher struct {
	Message string
}

// NewFailureMatcher returns a new *FailureMatcher.
func NewFailureMatcher(path, expected, actual string) *FailureMatcher {
	return &FailureMatcher{
		Message: fmt.Sprintf("path = %s, expected = %s, actual = %s", path, expected, actual),
	}
}

// Match implements types.GomegaMatcher.
func (matcher *FailureMatcher) Match(actual interface{}) (success bool, err error) {
	return false, nil
}

// FailureMessage implements types.GomegaMatcher.
func (matcher *FailureMatcher) FailureMessage(actual interface{}) (message string) {
	return matcher.Message
}

// NegatedFailureMessage implements types.GomegaMatcher.
func (matcher *FailureMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Not %s", matcher.Message)
}
