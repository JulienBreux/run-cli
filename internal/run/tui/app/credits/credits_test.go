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

package credits

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := tview.NewApplication()
	onClose := func() {}

	c := New(app, onClose)

	assert.NotNil(t, c)
	assert.Equal(t, app, c.app)
	assert.NotNil(t, c.authorColors)
	assert.Contains(t, c.contributors, "@JulienBreux")
	assert.Contains(t, c.contributors, "@rdout2")

	// Verify author colors initialized
	for _, author := range Contributors {
		assert.Contains(t, c.authorColors, author)
		assert.NotEqual(t, tcell.ColorDefault, c.authorColors[author])
	}
}

func TestAnimationControl(t *testing.T) {
	app := tview.NewApplication()
	c := New(app, func() {})

	// Ensure channels are ready
	assert.NotNil(t, c.stop)

	// Start Animation
	assert.NotPanics(t, func() {
		c.StartAnimation()
	})

	// Wait a bit to let goroutine run
	time.Sleep(50 * time.Millisecond)

	// Stop Animation
	assert.NotPanics(t, func() {
		c.StopAnimation()
	})

	// Verify channel closed
	select {
	case <-c.stop:
		// Channel closed, good
	default:
		t.Error("Stop channel should be closed")
	}
}

func TestInputHandler(t *testing.T) {
	app := tview.NewApplication()
	called := false
	onClose := func() { called = true }

	c := New(app, onClose)
	handler := c.InputHandler()

	// ESC should trigger onClose
	handler(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone), nil)
	assert.True(t, called, "Escape should close credits")

	// Reset
	called = false
	c = New(app, onClose)
	handler = c.InputHandler()

	// Enter should trigger onClose
	handler(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone), nil)
	assert.True(t, called, "Enter should close credits")
}

func TestUpdate(t *testing.T) {
	app := tview.NewApplication()
	c := New(app, func() {})

	// Initial state
	initialScroll := c.scrollY
	initialParticles := len(c.particles)

	// Force update with simulated rect
	c.SetRect(0, 0, 100, 100)
	// Reset the lastUpdate timer to stable offset to ensure dt > 0 (as per memory and stability fix)
	c.lastUpdate = time.Now().Add(-100 * time.Millisecond)
	c.update()

	// Check if particles were spawned
	assert.Greater(t, len(c.particles), initialParticles, "Should spawn particles")

	// Check scroll movement (might be negative)
	c.lastUpdate = time.Now().Add(-100 * time.Millisecond)
	c.update()
	assert.Less(t, c.scrollY, initialScroll, "Text should scroll upwards (negative Y)")
}

func TestDraw(t *testing.T) {
	app := tview.NewApplication()
	c := New(app, func() {})

	screen := tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	assert.NoError(t, err)

	c.SetRect(0, 0, 80, 24)

	// Simulate one update to spawn particles
	c.update()

	// Draw
	assert.NotPanics(t, func() {
		c.Draw(screen)
	})
}
