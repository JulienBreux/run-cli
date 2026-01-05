package format_test

import (
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
