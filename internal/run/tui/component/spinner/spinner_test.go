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

package spinner

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	app := tview.NewApplication()
	s := New(app, 0)
	assert.NotNil(t, s)
	assert.NotNil(t, s.TextView)
	assert.Equal(t, app, s.app)
}

func TestStartStop(t *testing.T) {
	// Use SimulationScreen to allow app.Run() without terminal
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}

	app := tview.NewApplication()
	app.SetScreen(screen)

	// Run app in goroutine to process QueueUpdateDraw events
	go func() {
		_ = app.Run()
	}()

	// Ensure app stops at end of test
	defer app.Stop()

	s := New(app, 0)

	// Start
	s.Start("Loading...")

	s.mu.Lock()
	assert.NotNil(t, s.cancel, "Cancel function should be set after Start")
	s.mu.Unlock()

	// Wait a bit to let goroutine run and process updates
	time.Sleep(150 * time.Millisecond)

	// Stop
	s.Stop("Done")

	s.mu.Lock()
	assert.Nil(t, s.cancel, "Cancel function should be nil after Stop")
	s.mu.Unlock()
}

func TestSetContext(t *testing.T) {
	app := tview.NewApplication()
	s := New(app, 0)

	s.SetContext("loading something")
	s.mu.Lock()
	assert.Equal(t, "loading something", s.context)
	s.mu.Unlock()
}
