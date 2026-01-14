package logo_test

import (
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
)

func TestString(t *testing.T) {
	s := logo.String()

	if !strings.Contains(s, "___ _   _ _  _") {
		t.Errorf("logo.String() should contain the logo art")
	}
}

func TestNew(t *testing.T) {
	l := logo.New()

	if l == nil {
		t.Error("logo.New() should return a non-nil TextView")
	}
	
	// We can't easily extract text from TextView directly without drawing, 
	// but we can check if it was initialized (not crashing).
	// In a real TUI test we might use a screen simulation, but for unit test ensuring it returns expected type is often enough.
	// However, we can assert properties if we want.
}
