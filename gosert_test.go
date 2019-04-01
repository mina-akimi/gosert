package gosert

import (
	"testing"

	"github.com/mina-akimi/gosert/v2/matcher"
)

func TestNewMatcherFromFile(t *testing.T) {
	m := MustMatcher(
		NewMatcherFromFile(
			"test_asset/golden1.json",
			map[string]string{
				"VAR":       "value1_0",
				"TIMESTAMP": "2018-10-05T12:13:14.000Z",
			},
			matcher.JSONParserInstance),
	)
	act := `{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "2018-10-05T12:13:14.123Z"
				}
			}
		}`

	matched, err := m.Match(act)
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
}
