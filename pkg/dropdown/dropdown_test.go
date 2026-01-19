package dropdown_test

import (
	"testing"

	"github.com/JulienBreux/run-cli/pkg/dropdown"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	d := dropdown.New()
	assert.NotNil(t, d)
}

func TestDraw(t *testing.T) {
	d := dropdown.New()
	d.SetRect(0, 0, 20, 1) // Set a rect so it has size
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)
	
	d.Draw(screen)
	// No panic implies success for now
}
