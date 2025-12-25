package format

import "gopkg.in/yaml.v3"

// ToYAML returns value in YAML
func ToYAML(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
