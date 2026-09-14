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

package header_test

import (
	"strings"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testInfo := info.Info{
		Project: "test-project",
		Region:  "us-central1",
		User:    "test-user",
	}

	h := header.New(testInfo)

	assert.NotNil(t, h)
	assert.Equal(t, 3, h.GetItemCount())
}

func TestUpdateInfo(t *testing.T) {
	// Initialize the component first to set the global infoView variable
	initialInfo := info.Info{
		Project: "p1",
		Region:  "r1",
		User:    "u1",
	}
	_ = header.New(initialInfo)

	// Now update it
	newInfo := info.Info{
		Project: "p2",
		Region:  "r2",
		User:    "u2",
	}
	
	// This function modifies the global infoView. 
	// Since we can't inspect the text content easily without drawing, we just ensure it doesn't panic.
	assert.NotPanics(t, func() {
		header.UpdateInfo(newInfo)
	})
	
	// Note: Testing side effects on global variables is brittle in parallel tests, 
	// but acceptable here given the legacy code structure.
}

func TestHeader_ShortcutsOrder(t *testing.T) {
	testInfo := info.Info{
		Project: "test-project",
		Region:  "us-central1",
		User:    "test-user",
	}

	h := header.New(testInfo)
	assert.NotNil(t, h)
	assert.Equal(t, 3, h.GetItemCount())

	// Item 1 is the shortcuts Flex
	shortcutsFlex, ok := h.GetItem(1).(*tview.Flex)
	assert.True(t, ok)
	assert.Equal(t, 2, shortcutsFlex.GetItemCount())

	// Item 1 of shortcutsFlex is col2 which contains resource shortcuts
	col2, ok := shortcutsFlex.GetItem(1).(*tview.TextView)
	assert.True(t, ok)

	text := col2.GetText(true)
	assert.Contains(t, text, "<ctrl+n>")
	assert.Contains(t, text, "<ctrl+s>")

	idxInstance := strings.Index(text, "<ctrl+n>")
	idxService := strings.Index(text, "<ctrl+s>")
	assert.Less(t, idxInstance, idxService, "Instances (<ctrl+n>) should appear before Services (<ctrl+s>) in header shortcuts column")
}

func TestHeader_UpdateNotice(t *testing.T) {
	testInfo := info.Info{
		Project: "test-project",
		Region:  "us-central1",
		User:    "test-user",
	}

	h := header.New(testInfo)
	assert.NotNil(t, h)

	infoCol, ok := h.GetItem(0).(*tview.TextView)
	assert.True(t, ok)

	text := infoCol.GetText(true)
	assert.Contains(t, text, "Version:")
	assert.NotContains(t, text, "update:")

	header.SetUpdateAvailable("v9.9.9")

	textWithUpdate := infoCol.GetText(true)
	assert.Contains(t, textWithUpdate, "Version:")
	assert.Contains(t, textWithUpdate, "(update: v9.9.9)")

	// Reset for subsequent tests
	header.SetUpdateAvailable("")
}

func TestStartUpdateCheck(t *testing.T) {
	testInfo := info.Info{
		Project: "test-project",
		Region:  "us-central1",
		User:    "test-user",
	}
	_ = header.New(testInfo)

	assert.NotPanics(t, func() {
		header.StartUpdateCheck(nil)
	})
}


