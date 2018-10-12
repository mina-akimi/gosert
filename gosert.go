package gosert

import (
	"fmt"
	"io/ioutil"

	"github.com/mina-akimi/gosert/matcher"
	"github.com/onsi/gomega/types"
)

// Matcher implements types.GomegaMatcher
type Matcher struct {
	parser   matcher.Parser
	expected matcher.Node
	// a current matcher that we delegate failure message to
	curMatcher types.GomegaMatcher
}

// NewMatcher returns a new matcher.  vars is used to replace variables in data.
func NewMatcher(data []byte, vars map[string]string, parser matcher.Parser) (*Matcher, error) {
	data, err := matcher.Replace(data, vars)
	if err != nil {
		return nil, err
	}

	return &Matcher{
		expected: matcher.Node{
			Type:  matcher.Object,
			Value: data,
		},
		parser: parser,
	}, err
}

// NewMatcher returns a new matcher.  The expected value is read from path.  vars is used to replace variables in data.
func NewMatcherFromFile(path string, vars map[string]string, parser matcher.Parser) (*Matcher, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return NewMatcher(bs, vars, parser)
}

// NewJSONMatcher returns a new matcher.
func NewJSONMatcher(data []byte, vars map[string]string) (*Matcher, error) {
	return NewMatcher(data, vars, matcher.JSONParserInstance)
}

// NewJSONMatcherFromFile returns a new matcher.
func NewJSONMatcherFromFile(path string, vars map[string]string) (*Matcher, error) {
	return NewMatcherFromFile(path, vars, matcher.JSONParserInstance)
}

// MustMatcher can be used with create matcher functions.  This panics if the create function returns err != nil.
func MustMatcher(m *Matcher, err error) *Matcher {
	if err != nil {
		panic(err)
	}
	return m
}

// Match matches a `*datatype.Timestamp`.
func (m *Matcher) Match(actual interface{}) (bool, error) {
	if actual == nil {
		return false, nil
	}

	var bs []byte
	bs, ok := actual.([]byte)
	if !ok {
		str, ok := actual.(string)
		if !ok {
			return false, fmt.Errorf("Unsupported actual type.  Only string or []byte.")
		}
		bs = []byte(str)
	}

	actNode := matcher.Node{
		Type:  matcher.Object,
		Value: bs,
	}

	mt, matched, err := matcher.Walk("", m.expected, actNode, m.parser)
	m.curMatcher = mt
	return matched, err
}

// FailureMessage returns failure message.
func (m *Matcher) FailureMessage(actual interface{}) string {
	return m.curMatcher.FailureMessage(actual)
}

// NegatedFailureMessage returns negated failure message.
func (m *Matcher) NegatedFailureMessage(actual interface{}) string {
	return m.curMatcher.NegatedFailureMessage(actual)
}
