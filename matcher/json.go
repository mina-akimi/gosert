package matcher

import (
	"encoding/json"
	"log"

	"github.com/buger/jsonparser"
)

// JSONParser is parser for JSON.
type JSONParser struct {
}

// ValidateObject implements Parser.
func (p *JSONParser) ValidateObject(data []byte) error {
	m := map[string]interface{}{}
	return json.Unmarshal(data, &m)
}

// GetFields implements Parser.
func (p *JSONParser) GetFields(data []byte) map[string]Node {
	m := map[string]Node{}
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		m[string(key)] = Node{
			Type:  jsonparserToInternal(dataType),
			Value: value,
		}
		return nil
	})
	return m
}

// GetArray implements Parser.
func (p *JSONParser) GetArray(data []byte) []Node {
	var nodes []Node
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			log.Printf("Error parsing array: %+v", err)
		}
		nodes = append(nodes, Node{
			Type:  jsonparserToInternal(dataType),
			Value: value,
		})
	})
	return nodes
}

// Delete implements Parser.
func (p *JSONParser) Delete(data []byte, key string) []byte {
	return jsonparser.Delete(data, key)
}

func jsonparserToInternal(dataType jsonparser.ValueType) ValueType {
	switch dataType {
	case jsonparser.NotExist:
		return NotExist
	case jsonparser.String:
		return String
	case jsonparser.Number:
		return Number
	case jsonparser.Object:
		return Object
	case jsonparser.Array:
		return Array
	case jsonparser.Boolean:
		return Boolean
	case jsonparser.Null:
		return Null
	default:
		return NotExist
	}
}
