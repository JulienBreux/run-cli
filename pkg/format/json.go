package format

import "encoding/json"

// ToJSON returns value in JSON
func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}
