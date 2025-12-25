package format

import "encoding/json"

// ToJSON returns value in JSON
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
