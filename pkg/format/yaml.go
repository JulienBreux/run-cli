package format

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// ToYAML returns value in YAML
func ToYAML(v any) ([]byte, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
