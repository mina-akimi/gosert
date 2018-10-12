<h1 align="center">
  <a href="https://github.com/mina-akimi/gosert"><img src="icons/gosert-logo-name.svg" alt="logo" width="60%" height="60%"/></a><br>
  <a href="https://github.com/mina-akimi/gosert">Gosert</a>
</h1>

<h4 align="center">A structured assertion library based on <a href="https://onsi.github.io/gomega/">gomega</a>.</h4>

The gosert library is an attempt to make assertions on structured objects nice and clean.  This project is influenced by the [golden file](https://medium.com/soon-london/testing-with-golden-files-in-go-7fccc71c43d3) idea.

Goals include:

* ***Simplicity***.  Dead simple API.

* ***Flexibility***.  Allow the user to express constraints on fields (not just hard-coded values).

* ***Presentation agnostic***.  Assume a minimal set of base data types, and use parser interface to work over the data (so JSON, YAML etc all work).

## Overview

The library allows the user to define golden files that looks largely the same as the actual data, with [gomega](https://onsi.github.io/gomega/) like DSL for defining constraints.

Example:

```
# Actual JSON
{
  "userId": "0001",
  "name": "Ethan Hunt",
  "age": 36,
  "interestes": ["climbing", "jogging", "marshal arts"],
  "address": ["IMF", "New York", "US"]
  "connections": [
    {
      "id": "0002",
      "name": "Julia Meade",
      "relationship": "ex"
    },
    {
      "id": "0003",
      "name": "Ilsa Faust",
      "relationship": "unclear"
    }
  ],
  "updatedAt": "2018-10-05T12:13:14.123Z"
}

# Golden File
{
  "userId": "0001",
  "name": "Ethan Hunt",
  "age": 36,
  "interestes": ["climbing", "marshal arts", "jogging"], // By default order doesn't matter
  "address": "{{Not(BeEmpty())}}", // Assert the array is not empty
  "connections": [
    {
      "_gst_id": "id=0002", // Match element containing "id": "0002"
      "id": "0002",
      "name": "Julia Meade",
      "relationship": "ex"
    },
    {
      "_gst_id": "id=0003",
      "id": "0003",
      "name": "Ilsa Faust",
      "relationship": "unclear"
    }
  ],
  "updatedAt": "{{BeTimestamp(2018-10-05T12:13:14.000Z, 5000)}}" // This means 2018-10-05T12:13:14.000Z +/- 5 seconds
}
```

Code example:

```
import (
  . "github.com/onsi/gomega"
  . "github.com/mina-akimi/gosert"
)

var _ = It("Test", func() {
    actual := `
    {
      "userId": "0001",
      "name": "Ethan Hunt",
      "age": 36,
      "interestes": ["climbing", "jogging", "marshal arts"],
      "address": ["IMF", "New York", "US"]
      "connections": [
        {
          "id": "0002",
          "name": "Julia Meade",
          "relationship": "ex"
        },
        {
          "id": "0003",
          "name": "Ilsa Faust",
          "relationship": "unclear"
        }
      ],
      "updatedAt": "2018-10-05T12:13:14.123Z"
    }
    `
    Expect(actual).To(MustMatcher(NewJSONMatcherFromFile("path/to/file", nil)))
})

```


## Data Types

The following base types are supported

* `String`
* `Number` (including integers, floats)
* `Boolean`
* `Object`
* `Array`

## Functions

Functions are string values enclosed in `{{}}`.

| Function                                    | Applicable Data Type | Meaning                                                                                                                       | Example                                              |
|---------------------------------------------|----------------------|-------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| `{{BeEmpty()}}`                             | `String`, `Array`    | The object is empty or the key is not present.                                                                                | These will all match: <br/> * "" <br/> * [] <br/> * The key is missing |
| `{{Not(BeEmpty())}}`                        | `String`, `Array`    | The object is not empty                                                                                                       |                                                      |
| `{{BeNumerically(<comparator>, <values>)}}` | `Number`             | See [here](https://onsi.github.io/gomega/#benumericallycomparator-string-compareto-interface)                                 | `{{BeNumerically(~, 123, 0.01)}}`                    |
| `{{BeTimestamp(<time>, <delta>)}}`          | `String`             | * `<time>` must be of [RFC3339 format](https://gobyexample.com/time-formatting-parsing) <br/>* `<delta>` is number of milliseconds | `{{BeTimestamp(2018-10-05T12:13:14.000Z, 5000)}}`    |

## Advanced Usage

### Variable Substitution

Variables can be defined in golden files in the form `${{MY_VAR}}`.  When creating a matcher, the user must supply a `map[string]string` mapping the variable names to values.

If any variable is not defined in the value map, an error is returned.

### Array Assertion

By default, ordering is ignored.

For base type array, the expected value can be an array of base values.

For object array, the expected values must contain a field called `_gst_id` with value `<key_name>=<value>` (see example above).

One requirement is that all expected objects in the same array must define the same `<key_name>`, otherwise an error occurs.

### Multipart File

You can define multiple objects (both fixture and expected objects) in a single file.

Example golden file:

```
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
```

Note each object must be preceded by a line like this:
```
### key=<object_name>, <my comments>
```

The the object can be read using `MultipartReader`:
```
r := MustReader(NewMultipartReaderFromFile(
    file,
    map[string]string{
        "MY_VAR": "baz",
        "NOW":    "2018-10-05T12:13:14.123Z",
    },
    matcher.JSONParserInstance,
))

Expect(r.GetData("my_fixture")).To(r.MustGetMatcher("my_matcher"))
```

## Examples

See https://github.com/mina-akimi/gosert/blob/master/matcher/dsl_test.go for more examples.