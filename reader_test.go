package gosert

import (
	"testing"

	"github.com/mina-akimi/gosert/matcher"
)

func TestMultipartReader_MustGetMatcher(t *testing.T) {
	data := []byte(`
		### key=my_fixture, my fixture
		{
			"foo": "bar",
			"baz": "${{MY_VAR}}",
			"quux": "${{NOW}}"
		}

		### key=my_matcher, my awesome matcher
		{
			"foo": "bar",
			"baz": "{{Not(BeEmpty())}}",
			"quux": "{{BeTimestamp(${{NOW}}, 5000)}}"
		}
	`)

	r := MustReader(NewMultipartReader(
		data,
		map[string]string{
			"MY_VAR": "baz",
			"NOW":    "2018-10-05T12:13:14.123Z",
		},
		matcher.JSONParserInstance,
	))

	matched, err := r.MustGetMatcher("my_matcher").Match(r.GetData("my_fixture"))
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
}

func TestMultipartReader_UpdateVars(t *testing.T) {
	data := []byte(`
		### key=my_fixture, my fixture
		{
			"foo": "bar",
			"baz": "${{MY_VAR}}",
			"quux": "${{NOW}}"
		}

		### key=my_matcher, my awesome matcher
		{
			"foo": "bar",
			"baz": "qux",
			"quux": "{{BeTimestamp(${{NOW}}, 5000)}}"
		}
	`)

	r := MustReader(NewMultipartReader(
		data,
		map[string]string{
			"MY_VAR": "baz",
			"NOW":    "2018-10-05T12:13:14.123Z",
		},
		matcher.JSONParserInstance,
	))

	matched, err := r.MustGetMatcher("my_matcher").Match(r.GetData("my_fixture"))
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}

	err = r.UpdateVars(map[string]string{
		"MY_VAR": "qux",
		"NOW":    "2018-10-05T12:13:14.123Z",
	})
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}

	matched, err = r.MustGetMatcher("my_matcher").Match(r.GetData("my_fixture"))
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
}
