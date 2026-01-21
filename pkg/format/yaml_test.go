/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package format_test

import (
	"fmt"
	"testing"

	"github.com/JulienBreux/run-cli/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestToYAML(t *testing.T) {
	data := struct {
		Version string `yaml:"version"`
	}{
		Version: "1.0.0",
	}
	actual, err := format.ToYAML(data)
	expected := "version: 1.0.0\n"

	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestToYAMLError(t *testing.T) {
	// YAML encoder fails on invalid map keys like functions or slices?
	// Or maybe just a channel.
	// yaml.v3 usually returns error for channels.
	data := make(chan int)
	_, err := format.ToYAML(data)
	assert.Error(t, err)
}

type FailMarshaler struct{}

func (f FailMarshaler) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("expected error")
}

func TestToYAML_MarshalFail(t *testing.T) {
	_, err := format.ToYAML(FailMarshaler{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected error")
}
