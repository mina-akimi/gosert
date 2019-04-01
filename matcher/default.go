package matcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// KeyID is used to identify elements by ID in an array
	KeyID = "_gst_id"
	// KeyIndex is used to identify elements by index in an array
	KeyIndex = "_gst_index"
)

var (
	// JSONParserInstance is a singleton.
	JSONParserInstance = &JSONParser{}
)

// ValueType defines value types available.
type ValueType int

// Nodes represents []Node
type Nodes []Node

const (
	NotExist = ValueType(iota)
	String
	Number
	Object
	Array
	Boolean
	Null
)

// String returns string.
func (t ValueType) String() string {
	switch t {
	case NotExist:
		return "NotExist"
	case String:
		return "String"
	case Number:
		return "Number"
	case Object:
		return "Object"
	case Array:
		return "Array"
	case Boolean:
		return "Boolean"
	case Null:
		return "Null"
	default:
		return "NotExist"
	}
}

// Node represents an abstract node in tree structure.
type Node struct {
	Type  ValueType
	Value []byte
}

// String returns string.
func (n Node) String() string {
	return fmt.Sprintf("type = %s, value = %s", n.Type.String(), string(n.Value))
}

// String returns string.
func (ns Nodes) String() string {
	var strs []string
	for _, node := range ns {
		strs = append(strs, node.String())
	}
	s := strings.Join(strs, "|")
	return fmt.Sprintf("[%s]", s)
}

// Parser represents a parser for data types.
type Parser interface {
	ValidateObject(data []byte) error
	GetFields(data []byte) map[string]Node
	GetArray(data []byte) []Node
	Delete(data []byte, key string) []byte
}

// ParseTime parses str with format time.RFC3339Nano.
func ParseTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, str)
}

// IsObjects returns true if all elements in nodes has type Object.
func IsObjects(nodes []Node) bool {
	for _, node := range nodes {
		if node.Type != Object {
			return false
		}
	}
	return true
}

// IsBaseType returns true if node has base type.
func IsBaseType(node Node) bool {
	switch node.Type {
	case String, Number, Boolean, Null:
		return true
	}
	return false
}

// IsObjects returns true if all elements in nodes has base type.
func IsBaseTypes(nodes []Node) bool {
	for _, node := range nodes {
		if !IsBaseType(node) {
			return false
		}
	}
	return true
}

func toNumber(input []byte) (float64, error) {
	s := string(input)
	return strconv.ParseFloat(s, 64)
}

func toBool(input []byte) (bool, error) {
	s := string(input)
	if strings.ToLower(s) == "true" {
		return true, nil
	}
	if strings.ToLower(s) == "false" {
		return false, nil
	}
	return false, fmt.Errorf("expected bool but got %s", s)
}
