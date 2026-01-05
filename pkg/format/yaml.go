package format

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// ToYAML returns value in YAML
func ToYAML(v any) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("yaml encode panic: %v", r)
		}
	}()

	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
