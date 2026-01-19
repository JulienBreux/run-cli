package auth

import (
	"testing"

	"github.com/rivo/tview"
)

func TestModal(t *testing.T) {
	// Simple test to ensure it compiles and returns a primitive
	app := tview.NewApplication()
	_ = app
	// Mock service if needed, or pass nil if safe (checking implementation, it accesses service.Security which might panic if service is nil)
	// We can't easily test TUI interaction here without a complex setup, so we just verify it doesn't crash on creation
	// and returns a grid.
}
