package matcher

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

var (
	// Usage: {{ts(2018-01-02T12:13:14.123Z, 5000)}}, which means 2018-01-02T12:13:14.123Z +/- 5000 milliseconds
	patternTimestamp = regexp.MustCompile(`{{BeTimestamp\((?P<time>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z),\s*(?P<delta>\d+)\)}}`)
	// Usage: {{num(123.23, 0.5)}}, which means 123.23 +/- 0.5
	patternNumber = regexp.MustCompile(`{{BeNumerically\((?P<comparator>[^,]+),\s*(?P<data>[^)]+)\)}}`)
	// Usage: {{BeEmpty()}}, which means the string/array must be empty
	patternEmpty = regexp.MustCompile(`{{BeEmpty\(\)}}`)
	// Usage: {{Not(BeEmpty())}}, which means the string/array must not be empty
	patternNotEmpty = regexp.MustCompile(`{{Not\(BeEmpty\(\)\)}}`)
	// Usage: ${{MY_VAR}}, which can be replaced with a value
	patternVariable = regexp.MustCompile(`\${{(?P<var>\w+)}}`)

	arrayPatterns = []string{"{{Not(BeEmpty())}}"}
)

// ===========
// Tree walker
// ===========

// Walk recursively iterates the tree structure, matching elements in act with exp.
func Walk(path string, exp, act Node, parser Parser) (types.GomegaMatcher, bool, error) {
	switch act.Type {
	case String:
		if exp.Type != String {
			return NewFailureMatcher(path, exp.Type.String(), act.Type.String()), false, fmt.Errorf("path has type String but assertion uses %s", exp.Type.String())
		}
		matcher, err := CreateStringMatcher(string(exp.Value))
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		matched, err := matcher.Match(string(act.Value))
		if !matched || err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), matched, err
		}
	case Number:
		actVal, err := toNumber(act.Value)
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		switch exp.Type {
		case String:
			matcher, err := CreateNumberMatcher(string(exp.Value))
			if err != nil {
				return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
			}
			matched, err := matcher.Match(actVal)
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), matched, err
		case Number:
			expVal, err := toNumber(exp.Value)
			if err != nil {
				return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
			}
			matcher := &matchers.BeNumericallyMatcher{
				Comparator: "~",
				CompareTo:  []interface{}{expVal, 0.05},
			}
			matched, err := matcher.Match(actVal)
			if !matched || err != nil {
				return NewFailureMatcher(path, string(exp.Value), string(act.Value)), matched, err
			}
		default:
			return NewFailureMatcher(path, exp.Type.String(), act.Type.String()), false, fmt.Errorf("path has type Number but assertion uses %s.  Allowed types are String or Number.", exp.Type.String())
		}
	case Boolean:
		if exp.Type != Boolean {
			return NewFailureMatcher(path, exp.Type.String(), act.Type.String()), false, fmt.Errorf("path has type Boolean but assertion uses %s", exp.Type.String())
		}
		expBool, err := toBool(exp.Value)
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		actBool, err := toBool(act.Value)
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		if expBool != actBool {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
	case Array:
		switch exp.Type {
		case Array:
			return MatchArrayWithArray(path, parser.GetArray(exp.Value), parser.GetArray(act.Value), parser)
		case String:
			return MatchArrayWithString(path, exp, parser.GetArray(act.Value))
		default:
			return NewFailureMatcher(path, exp.Type.String(), act.Type.String()), false, fmt.Errorf("unsupported expected value type '%s' for Array type.  See Gosert doc.", exp.Type.String())
		}
	case Object:
		err := parser.ValidateObject(act.Value)
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		err = parser.ValidateObject(exp.Value)
		if err != nil {
			return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, err
		}
		actObj := parser.GetFields(act.Value)
		expObj := parser.GetFields(exp.Value)
		for k, v := range expObj {
			if a, ok := actObj[k]; !ok {
				if isBeEmpty(v) {
					continue
				}
				return NewFailureMatcher(path, string(exp.Value), string(act.Value)), false, nil
			} else {
				// Recursion
				matcher, matched, err := Walk(path+"."+k, v, a, parser)
				if !matched || err != nil {
					return matcher, matched, err
				}
			}
		}
	}
	return SuccessMatcherInstance, true, nil
}

// MatchArrayWithArray matches act with exp as plain arrays.
func MatchArrayWithArray(path string, exp, act []Node, parser Parser) (types.GomegaMatcher, bool, error) {
	if IsBaseTypes(exp) {
		if !IsBaseTypes(act) {
			return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, fmt.Errorf("array should contain base type only but got object type")
		}
		if !baseNodesEqual(exp, act) {
			return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, nil
		}
		return SuccessMatcherInstance, true, nil
	} else if IsObjects(exp) {
		metaKey, expMap, err := createExpectedIDMapper(exp, parser)
		if err != nil {
			return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, err
		}
		actMap, err := createActualIDMapper(act, metaKey, parser)
		if err != nil {
			return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, err
		}
		for k, e := range expMap {
			a, ok := actMap[k]
			if !ok {
				return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, nil
			}
			// Recursion
			matcher, matched, err := Walk(path, e, a, parser)
			if !matched || err != nil {
				return matcher, matched, err
			}
		}
		return SuccessMatcherInstance, true, nil
	} else {
		return NewFailureMatcher(path, Nodes(exp).String(), Nodes(act).String()), false, fmt.Errorf("array type assertion must all be base data type or all object type, not a mixtureof both")
	}
}

// MatchArrayWithArray matches act with exp.  exp must be a function.
func MatchArrayWithString(path string, exp Node, act []Node) (types.GomegaMatcher, bool, error) {
	if exp.Type != String {
		return NewFailureMatcher(path, exp.String(), Nodes(act).String()), false, fmt.Errorf("array type assertion cannot use type '%s' as expected value", exp.Type.String())
	}
	v := string(exp.Value)
	if patternEmpty.MatchString(v) {
		if len(act) > 0 {
			return NewFailureMatcher(path, exp.String(), Nodes(act).String()), false, nil
		}
		return SuccessMatcherInstance, true, nil
	}
	if patternNotEmpty.MatchString(v) {
		if len(act) <= 0 {
			return NewFailureMatcher(path, exp.String(), Nodes(act).String()), false, nil
		}
		return SuccessMatcherInstance, true, nil
	}
	return NewFailureMatcher(path, exp.String(), Nodes(act).String()), false, fmt.Errorf("array type assertion can only use functions %+v", arrayPatterns)
}

// CreateNumberMatcher returns a matcher for numbers.  input must be a function.
func CreateNumberMatcher(input string) (types.GomegaMatcher, error) {
	if patternNumber.MatchString(input) {
		m := patternNumber.FindStringSubmatch(input)
		sub := mapSubexpNames(m, patternNumber.SubexpNames())
		comparator := sub["comparator"]
		data := sub["data"]
		datas := strings.Split(data, ",")
		var flts []interface{}
		for _, d := range datas {
			d = strings.TrimSpace(d)
			dflt, err := strconv.ParseFloat(d, 64)
			if err != nil {
				return nil, err
			}
			flts = append(flts, dflt)
		}
		return &matchers.BeNumericallyMatcher{
			Comparator: comparator,
			CompareTo:  flts,
		}, nil
	}
	return nil, fmt.Errorf("path has type Number but assertion ('%s') does not have the correct format.  Must be {{BeNumerically(...)}}.  See Gosert doc.", input)
}

// CreateStringMatcher returns a matcher for string.  input must be either a plain string or a function.
func CreateStringMatcher(input string) (types.GomegaMatcher, error) {
	if patternEmpty.MatchString(input) {
		return &matchers.BeEmptyMatcher{}, nil
	}
	if patternNotEmpty.MatchString(input) {
		return &matchers.NotMatcher{
			Matcher: &matchers.BeEmptyMatcher{},
		}, nil
	}
	if patternTimestamp.MatchString(input) {
		m := patternTimestamp.FindStringSubmatch(input)
		sub := mapSubexpNames(m, patternTimestamp.SubexpNames())
		tsStr := sub["time"]
		deltaStr := sub["delta"]
		ts, err := ParseTime(tsStr)
		if err != nil {
			return nil, err
		}
		delta, err := strconv.ParseInt(deltaStr, 10, 64)
		if err != nil {
			return nil, err
		}
		return NewTimestampMatcher(ts, time.Duration(delta)*time.Millisecond), nil
	}
	return &matchers.EqualMatcher{
		Expected: input,
	}, nil
}

func isBeEmpty(node Node) bool {
	return node.Type == String && patternEmpty.Match(node.Value)
}

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i, _ := range n {
		r[n[i]] = m[i]
	}
	return r
}

func baseNodesEqual(exp, act []Node) bool {
	return nodesContain(exp, act) && nodesContain(act, exp)
}

type nodeCounter struct {
	node Node
	seen bool
}

func nodesContain(set, subset []Node) bool {
	var counters []*nodeCounter
	for _, s := range set {
		counters = append(counters, &nodeCounter{
			node: s,
			seen: false,
		})
	}

	for _, s := range subset {
		var found bool
		for _, c := range counters {
			if s.Type == c.node.Type && bytes.Equal(s.Value, c.node.Value) && !c.seen {
				c.seen = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// createExpectedIDMapper takes a list of nodes and returns the ID key and a map from the key to Node.
//
// E.g., if we have a Node with value like this
//
//     {
//       "_gst_id: "orderId=1234",
//       "field0": "value0"
//     }
//
// Then key will be "orderId" and the map will contain "1234" mapped to this Node (with "_gts_id" field removed, for matching)
//
// Returns an error if the key is not the same for all Nodes, e.g., one Node has `"_gst_id": "orderId="1234"` and another has `"_gst_id": "jobId="2345"`.
func createExpectedIDMapper(nodes []Node, parser Parser) (string, map[string]Node, error) {
	result := map[string]Node{}
	var metaKey string
	for _, node := range nodes {
		m := parser.GetFields(node.Value)
		if name, ok := m[KeyID]; ok {
			if name.Type != String {
				return "", nil, fmt.Errorf("'%s' field must be of type String, was %s", KeyID, name.Type.String())
			}
			parts := strings.Split(string(name.Value), "=")
			if len(parts) != 2 {
				return "", nil, fmt.Errorf("'%s' field must have value of format 'key=value', was %s", KeyID, string(name.Value))
			}
			mKey := parts[0]
			if metaKey != "" && mKey != metaKey {
				return "", nil, fmt.Errorf("all elements in the same array must have the same '%s' key part", KeyID)
			}
			metaKey = mKey
			key := parts[1]
			// Must delete KeyID to avoid match failure later, since it's not part of actual value
			node.Value = parser.Delete(node.Value, KeyID)
			result[key] = node
		} else {
			return "", nil, fmt.Errorf("object array assertion must provide a key '%s' for each element.  See Gosert doc.", KeyID)
		}
	}
	return metaKey, result, nil
}

func createActualIDMapper(nodes []Node, key string, parser Parser) (map[string]Node, error) {
	result := map[string]Node{}
	for _, node := range nodes {
		m := parser.GetFields(node.Value)
		if name, ok := m[key]; ok {
			result[string(name.Value)] = node
		} else {
			return nil, fmt.Errorf("object array assertion must have a key '%s' for each element.  See Gosert doc.", key)
		}
	}
	return result, nil
}

// =====================
// Variable substitution
// =====================

// Replace returns data with all variables replaced with values in vars.  If any variable is not defined in vars, an error is returned.
func Replace(data []byte, vars map[string]string) ([]byte, error) {
	ms := patternVariable.FindAllSubmatch(data, -1)
	ns := map[string]bool{}
	for _, m := range ms {
		ns[string(m[1])] = true
	}
	var names []string
	for k, _ := range ns {
		names = append(names, k)
	}

	result := string(data)
	for _, name := range names {
		value, ok := vars[name]
		if !ok {
			return nil, fmt.Errorf("variable '%s' undefined in substitution", name)
		}
		result = strings.Replace(result, fmt.Sprintf("${{%s}}", name), value, -1)
	}

	return []byte(result), nil
}
