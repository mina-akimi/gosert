package matcher

import (
	"testing"
)

func TestWalk_String(t *testing.T) {
	exp := Node{
		Type:  String,
		Value: []byte("hello"),
	}
	act := Node{
		Type:  String,
		Value: []byte("hello"),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Empty(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "{{BeEmpty()}}",
				"field1_1": "{{BeEmpty()}}",
				"field1_2": "{{BeEmpty()}}"
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [],
				"field1_2": ""
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Object_Success(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_0"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_0"
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Object_Number(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "{{BeNumerically(~, 123, 0.5)}}"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": 123.123
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Object_Failure_Number(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "{{BeNumerically(~, 123, 0.01)}}"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": 123.123
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should not nil but was %+v", err)
	}
}

func TestWalk_Object_Timestamp(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "{{BeTimestamp(2018-10-05T12:13:14.000Z, 5000)}}"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "2018-10-05T12:13:14.123Z"
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Object_Failure_Timestamp(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "{{BeTimestamp(2018-10-05T12:13:05.000Z, 5000)}}"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "2018-10-05T12:13:14.123Z"
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should not nil but was %+v", err)
	}
}

func TestWalk_Object_Failure_FieldMismatch(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_1"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_0"
				}
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Object_Failure_FieldMissing(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_1"
				}
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_0"
				},
				"field1_2": "value1_2"
			},
			"field2": []
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Failure_BadJSON(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		very bad json
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "value1_0",
				"field1_1": {
					"field1_1_0": "value1_1_0"
				},
				"field1_2": "value1_2"
				somethingbad
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err == nil {
		t.Fatalf("err should not be nil but was %+v", err)
	}
}

func TestWalk_Array_Orderless_String(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "foo", "world"]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "world", "foo"]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Orderless_String_Duplicates(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "foo", "world", "foo"]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "world", "foo", "foo"]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Orderless_String_Fail(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "foo", "world", "foo"]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "world", "foo"]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Orderless_Number(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [123, 456, 789.123]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [456, 123, 789.123]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Object_ByID(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"_gst_id": "myId=id0",
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"_gst_id": "myId=id1",
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"_gst_id": "myId=id2",
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Object_ByIndex(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"_gst_index": 0,
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"_gst_index": 2,
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Object_Fail_NoGstID(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"_gst_id": "myId=id0",
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"_gst_id": "myId=id2",
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err == nil {
		t.Fatalf("err should not be nil but was %+v", err)
	}
}

func TestWalk_Array_Object_Fail_DifferentMetaKey(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"_gst_id": "myId=id0",
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"_gst_id": "myOtherId=id2",
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"_gst_id": "myId=id2",
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": [
					{
						"myId": "id0",
						"field1_0_0": 100,
						"field1_0_1": "hello",
						"field1_0_2": {
							"field1_0_2_0": "value1_0_2_0"
						}
					},
					{
						"myId": "id1",
						"field1_0_0": 101.123,
						"field1_0_1": "world",
						"field1_0_2": {
							"field1_0_2_1": "value1_0_2_1"
						}
					},
					{
						"myId": "id2",
						"field1_0_0": 102.999,
						"field1_0_1": "foo",
						"field1_0_2": {
							"field1_0_2_2": "value1_0_2_2"
						}
					}
				]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err == nil {
		t.Fatalf("err should not be nil but was %+v", err)
	}
}

func TestWalk_Array_NotEmpty(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "{{Not(BeEmpty())}}"
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": ["hello", "world", "foo", "foo"]
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher != SuccessMatcherInstance {
		t.Fatalf("matcher should be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if !matched {
		t.Fatalf("matched should be true")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestWalk_Array_Fail_NotEmpty(t *testing.T) {
	exp := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": "{{Not(BeEmpty())}}"
			}
		}
	`),
	}
	act := Node{
		Type: Object,
		Value: []byte(`
		{
			"field0": "value0",
			"field1": {
				"field1_0": []
			}
		}
	`),
	}
	matcher, matched, err := Walk("", exp, act, JSONParserInstance)
	if matcher == SuccessMatcherInstance {
		t.Fatalf("matcher should not be SuccessMatcherInstance but was %+v, matched = %t, err = %+v", matcher, matched, err)
	}
	if matched {
		t.Fatalf("matched should be false")
	}
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
}

func TestReplace(t *testing.T) {
	data := []byte(`Hello world ${{FOO}} ${{BAR}} {{BeTimestamp(${{TIMESTAMP}}, 5000)}}`)
	vars := map[string]string{
		"FOO":       "Ethan",
		"BAR":       "Hunt",
		"TIMESTAMP": "2018-10-05T12:13:14.000Z",
	}

	res, err := Replace(data, vars)
	if err != nil {
		t.Fatalf("err should be nil but was %+v", err)
	}
	if string(res) != `Hello world Ethan Hunt {{BeTimestamp(2018-10-05T12:13:14.000Z, 5000)}}` {
		t.Fatalf("Replace() result incorrect, %s", string(res))
	}
}

func TestReplace_Failure_NoValue(t *testing.T) {
	data := []byte(`Hello world ${{FOO}} ${{BAR}} {{BeTimestamp(${{TIMESTAMP}}, 5000)}}`)
	vars := map[string]string{
		"FOO": "Ethan",
		"BAR": "Hunt",
	}

	_, err := Replace(data, vars)
	if err == nil {
		t.Fatalf("err should not be nil but was %+v", err)
	}
}
