package gosert

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/mina-akimi/gosert/matcher"
)

var (
	patternCommentStart = "# "
	patternHeaderStart  = "###"
	patternHeaderKey    = regexp.MustCompile(`key=(?P<key>\w+)`)
)

// Read reads a file with variables replaced.
func Read(path string, vars map[string]string) ([]byte, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return matcher.Replace(bs, vars)
}

// MustRead panics if error occurs.
func MustRead(path string, vars map[string]string) []byte {
	bs, err := Read(path, vars)
	if err != nil {
		panic(err)
	}
	return bs
}

// MultipartReader represents a read for file containing multiple objects (some as fixtures and some as expected values)
//
// Example:
//
//     ### key=my_fixture, my fixture
//     {
//       "foo": "bar",
//       "baz": "${{MY_VAR}}",
//       "quux": "${{NOW}}"
//     }
//
//     # This is a comment.  Any line starting with `# ` is a comment and is ignored (note the space after hash).
//     ### key=my_matcher, my awesome matcher
//     {
//       "foo": "bar",
//       "baz": "{{Not(BeEmpty())}}",
//       "quux": "{{BeTimestamp(${{NOW}}, 5000)}}"
//     }
type MultipartReader struct {
	raw    []byte
	parts  map[string][]byte
	parser matcher.Parser
}

// NewMultipartReader returns a new reader.
func NewMultipartReader(data []byte, vars map[string]string, parser matcher.Parser) (*MultipartReader, error) {
	reader := bytes.NewReader(data)
	return newMultipartReader(reader, vars, parser)
}

// NewMultipartReader returns a new reader.
func NewMultipartReaderFromFile(path string, vars map[string]string, parser matcher.Parser) (*MultipartReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return newMultipartReader(file, vars, parser)
}

// MustReader panics if an error occurs.
func MustReader(r *MultipartReader, err error) *MultipartReader {
	if err != nil {
		panic(err)
	}
	return r
}

func newMultipartReader(reader io.Reader, vars map[string]string, parser matcher.Parser) (*MultipartReader, error) {
	scanner := bufio.NewScanner(reader)
	var key string
	var object []byte
	var raw []byte
	parts := map[string][]byte{}
	for scanner.Scan() {
		raw = append(raw, scanner.Bytes()...)
		raw = append(raw, []byte(fmt.Sprintln())...)
		s := strings.TrimSpace(scanner.Text())
		if s == "" || strings.HasPrefix(s, patternCommentStart) {
			continue
		}
		if strings.HasPrefix(s, patternHeaderStart) {
			// Add to parts
			if len(object) > 0 {
				if key == "" {
					return nil, fmt.Errorf("multipart file cannot have header with empty body.  See Gosert doc.")
				}
				replaced, err := matcher.Replace(object, vars)
				if err != nil {
					return nil, err
				}
				parts[key] = replaced
			}

			// New key and clear lines
			m := patternHeaderKey.FindStringSubmatch(s)
			if len(m) > 0 {
				key = m[1]
			}
			object = nil
		} else {
			if key == "" {
				return nil, fmt.Errorf("section multipart file must have key.  See Gosert doc.")
			}
			object = append(object, scanner.Bytes()...)
			object = append(object, []byte(fmt.Sprintln())...)
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	if len(object) > 0 {
		bs, err := matcher.Replace(object, vars)
		if err != nil {
			return nil, err
		}
		parts[key] = bs
	}

	return &MultipartReader{
		raw:    raw,
		parts:  parts,
		parser: parser,
	}, nil
}

// GetData returns data with variables substituted.
func (r *MultipartReader) GetData(key string) []byte {
	return r.parts[key]
}

// GetData returns *Matcher with variables substituted.
func (r *MultipartReader) GetMatcher(key string) (*Matcher, error) {
	if bs, ok := r.parts[key]; ok {
		return NewMatcher(bs, nil, r.parser)
	}
	return nil, fmt.Errorf("no such key '%s' in file.  See Gosert doc.", key)
}

// MustGetMatcher panics if an error occurs.
func (r *MultipartReader) MustGetMatcher(key string) *Matcher {
	m, err := r.GetMatcher(key)
	if err != nil {
		panic(err)
	}
	return m
}

// UpdateVars updates r's variable substitution.  New matchers must be generated to take effect.
func (r *MultipartReader) UpdateVars(vars map[string]string) error {
	nr, err := newMultipartReader(bytes.NewReader(r.raw), vars, r.parser)
	if err != nil {
		return err
	}
	r.parts = nr.parts
	return nil
}
