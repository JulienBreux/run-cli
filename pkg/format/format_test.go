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
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/JulienBreux/run-cli/pkg/format"
	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	// Preparing test
	w := &bytes.Buffer{}
	v := struct {
		Message string `yaml:"message" json:"message"`
	}{
		Message: "Hello Hut!",
	}
	var c format.Callback = func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%s", "I'm just the callback!")
	}

	// Test JSON
	format.Print(w, format.JSON, v, c)
	assert.Equal(t, w.String(), "{\"message\":\"Hello Hut!\"}")
	w.Reset()

	// Test YAML
	format.Print(w, format.YAML, v, c)
	assert.Equal(t, w.String(), "message: Hello Hut!\n")
	w.Reset()

	// Test UNKNOWN
	format.Print(w, format.CUSTOM, v, c)
	assert.Equal(t, w.String(), "I'm just the callback!")
	w.Reset()
}

func TestPrintError(t *testing.T) {
	w := &bytes.Buffer{}
	invalidData := make(chan int)

	// Test JSON Error (silenced)
	format.Print(w, format.JSON, invalidData, nil)
	assert.Empty(t, w.String())

	// Test YAML Error (silenced)
	format.Print(w, format.YAML, invalidData, nil)
	assert.Empty(t, w.String())
}

func TestStringToFormat(t *testing.T) {
	assert.Equal(t, format.StringToFormat("yaml"), format.YAML)
	assert.Equal(t, format.StringToFormat("json"), format.JSON)
	assert.Equal(t, format.StringToFormat("breux"), format.CUSTOM)
}
