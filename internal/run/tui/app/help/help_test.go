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

package help

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestHelpModal(t *testing.T) {
	app := tview.NewApplication()
	closeFuncCalled := false
	closeFunc := func() {
		closeFuncCalled = true
	}

	modal := HelpModal(app, closeFunc)
	assert.NotNil(t, modal)
	assert.NotNil(t, modal.Grid)
	assert.NotNil(t, modal.Table)

	// Test Close handler
	handler := modal.Table.GetInputCapture()
	assert.NotNil(t, handler)

	// Escape
	handler(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone))
	assert.True(t, closeFuncCalled)

	// Enter
	closeFuncCalled = false
	handler(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
	assert.True(t, closeFuncCalled)

	// ?
	closeFuncCalled = false
	handler(tcell.NewEventKey(tcell.KeyRune, '?', tcell.ModNone))
	assert.True(t, closeFuncCalled)
}