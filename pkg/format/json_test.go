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
	"testing"

	"github.com/JulienBreux/run-cli/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	data := struct {
		Version string `json:"version"`
	}{
		Version: "1.0.0",
	}
	actual, err := format.ToJSON(data)
	expected := "{\"version\":\"1.0.0\"}"

	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestToJSONError(t *testing.T) {
	data := make(chan int)
	_, err := format.ToJSON(data)
	assert.Error(t, err)
}
