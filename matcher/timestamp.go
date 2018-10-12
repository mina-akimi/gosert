package matcher

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

// TimestampMatcher matches timestamps.
type TimestampMatcher struct {
	matcher matchers.BeTemporallyMatcher
}

// NewTimestampMatcher returns a matcher that checks datatype.Timestamp and `expected` is within `threshold`.
func NewTimestampMatcher(expected time.Time, threshold time.Duration) types.GomegaMatcher {
	return &TimestampMatcher{
		matcher: matchers.BeTemporallyMatcher{
			Comparator: "~",
			CompareTo:  expected,
			Threshold:  []time.Duration{threshold},
		},
	}
}

// Match matches a `string` formatted with RFC3339Nano.
func (matcher *TimestampMatcher) Match(actual interface{}) (bool, error) {
	if actual == nil {
		return false, nil
	}
	str, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("TimestampMatcher expects a string")
	}
	timestamp, err := ParseTime(str)
	if err != nil {
		return false, err
	}
	return matcher.matcher.Match(timestamp)
}

// FailureMessage returns failure message.
func (matcher *TimestampMatcher) FailureMessage(actual interface{}) string {
	return matcher.matcher.FailureMessage(actual)
}

// NegatedFailureMessage returns negated failure message.
func (matcher *TimestampMatcher) NegatedFailureMessage(actual interface{}) string {
	return matcher.matcher.NegatedFailureMessage(actual)
}
