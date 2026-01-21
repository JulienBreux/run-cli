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

package describe

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDescribeModal(t *testing.T) {
	app := tview.NewApplication()
	resource := map[string]string{"key": "value"}
	closeFunc := func() {}

	describer := DescribeModal(app, resource, "My Resource", closeFunc)

	assert.NotNil(t, describer)
	assert.NotNil(t, describer.TextView)
	assert.NotNil(t, describer.Content)
	
	// Check Primitive interface compliance
	var _ tview.Primitive = describer
}

func TestDescriber_Content(t *testing.T) {
	app := tview.NewApplication()
	resource := map[string]string{"foo": "bar"}
	
	describer := DescribeModal(app, resource, "Test", func(){})
	
	text := describer.TextView.GetText(true)
	// YAML output should contain "foo: bar"
	assert.Contains(t, text, "foo")
	assert.Contains(t, text, "bar")
}

func TestDescriber_InputCapture(t *testing.T) {
	app := tview.NewApplication()
	closed := false
	closeFunc := func() { closed = true }
	
	describer := DescribeModal(app, "data", "title", closeFunc)
	
	handler := describer.Content.GetInputCapture()
	assert.NotNil(t, handler)
	
	// Test Escape
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	ret := handler(eventEsc)
	assert.Nil(t, ret)
	assert.True(t, closed)
	
	// Reset
	closed = false
	
	// Test 'q'
	eventQ := tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone)
	ret = handler(eventQ)
	assert.Nil(t, ret)
	assert.True(t, closed)
	
	// Test other key
	eventOther := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	ret = handler(eventOther)
	assert.Equal(t, eventOther, ret)
	assert.True(t, closed) // Should be true because 'a' doesn't close it, wait...
	// Ah, I need to reset 'closed' logic if I want to assert it stays false.
}

func TestDescriber_InputCapture_NonClosing(t *testing.T) {
	app := tview.NewApplication()
	closed := false
	closeFunc := func() { closed = true }
	
	describer := DescribeModal(app, "data", "title", closeFunc)
	handler := describer.Content.GetInputCapture()
	
	eventOther := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	ret := handler(eventOther)
	assert.Equal(t, eventOther, ret)
	assert.False(t, closed)
}