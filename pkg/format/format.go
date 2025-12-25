package format

import (
	"fmt"
	"io"
)

// Format represents a format.
type Format int

// Callback represents a callback.
type Callback func(w io.Writer)

const (
	// JSON represents the JSON format.
	JSON Format = iota + 1
	// YAML represents the YAML format.
	YAML
	// CUSTOM represents a custom format and use callback.
	CUSTOM
)

// Print prints a formatted values.
func Print(w io.Writer, f Format, v any, c Callback) {
	switch f {
	case JSON:
		if b, err := ToJSON(v); err == nil {
			_, _ = fmt.Fprint(w, string(b))
		}
	case YAML:
		if b, err := ToYAML(v); err == nil {
			_, _ = fmt.Fprint(w, string(b))
		}
	case CUSTOM:
		c(w)
	}
}

// StringToFormat converts string format to typed format.
func StringToFormat(f string) Format {
	switch f {
	case "json":
		return JSON
	case "yaml":
		return YAML
	default:
		return CUSTOM
	}
}
