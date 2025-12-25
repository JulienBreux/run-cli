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

func TestStringToFormat(t *testing.T) {
	assert.Equal(t, format.StringToFormat("yaml"), format.YAML)
	assert.Equal(t, format.StringToFormat("json"), format.JSON)
	assert.Equal(t, format.StringToFormat("breux"), format.CUSTOM)
}
