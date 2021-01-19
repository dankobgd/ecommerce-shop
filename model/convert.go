package model

import (
	"encoding/json"
	"io"

	"github.com/gorilla/schema"
)

// SchemaDecoder to decode maps to struct
var SchemaDecoder = schema.NewDecoder()

// MapStrStrToJSON converts a map to a json string
func MapStrStrToJSON(obj map[string]string) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

// MapStrStrFromJSON will decode the key/value pair map
func MapStrStrFromJSON(data io.Reader) map[string]string {
	decoder := json.NewDecoder(data)

	var obj map[string]string
	if err := decoder.Decode(&obj); err != nil {
		return make(map[string]string)
	}
	return obj
}

// MapStrBoolToJSON converts a map to a json string
func MapStrBoolToJSON(obj map[string]bool) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

// MapStrBoolFromJSON will decode the key/value pair map
func MapStrBoolFromJSON(data io.Reader) map[string]bool {
	decoder := json.NewDecoder(data)

	var obj map[string]bool
	if err := decoder.Decode(&obj); err != nil {
		return make(map[string]bool)
	}
	return obj
}

// StrArrayToJSON converts an array to a json string
func StrArrayToJSON(obj []string) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

// StrSliceFromJSON will decode the array
func StrSliceFromJSON(data io.Reader) []string {
	decoder := json.NewDecoder(data)

	var obj []string
	if err := decoder.Decode(&obj); err != nil {
		return make([]string, 0)
	}
	return obj
}

// IntSliceFromJSON will decode the array
func IntSliceFromJSON(data io.Reader) []int {
	decoder := json.NewDecoder(data)

	var obj []int
	if err := decoder.Decode(&obj); err != nil {
		return make([]int, 0)
	}
	return obj
}

// MapStrInterfaceToJSON converts the map to a json string
func MapStrInterfaceToJSON(obj map[string]interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

// StrArrayFromInterface will decode the array
func StrArrayFromInterface(data interface{}) []string {
	stringArray := []string{}

	dataArray, ok := data.([]interface{})
	if !ok {
		return stringArray
	}

	for _, v := range dataArray {
		if str, ok := v.(string); ok {
			stringArray = append(stringArray, str)
		}
	}

	return stringArray
}

// MapStrInterfaceFromJSON will decode the key/value pair map
func MapStrInterfaceFromJSON(data io.Reader) map[string]interface{} {
	decoder := json.NewDecoder(data)

	var obj map[string]interface{}
	if err := decoder.Decode(&obj); err != nil {
		return make(map[string]interface{})
	}
	return obj
}

// StrToJSON converts a string to a json string
func StrToJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// StrFromJSON decodes the json string
func StrFromJSON(data io.Reader) string {
	decoder := json.NewDecoder(data)

	var s string
	if err := decoder.Decode(&s); err != nil {
		return ""
	}
	return s
}
