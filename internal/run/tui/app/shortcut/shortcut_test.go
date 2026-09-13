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

package shortcut

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	assert.NotEmpty(t, Registry)

	// Check for a few expected shortcuts
	foundHelp := false
	foundInstance := false
	for _, s := range Registry {
		if s.Key == "?" && s.Rune == '?' {
			foundHelp = true
		}
		if s.Key == "ctrl+n" && s.TCellKey == tcell.KeyCtrlN && s.Description == "Instances" && s.Category == CategoryGlobal {
			foundInstance = true
		}
	}
	assert.True(t, foundHelp, "Registry should contain the help shortcut '?'")
	assert.True(t, foundInstance, "Registry should contain 'ctrl+n' for Instances")
}

func TestRegistry_GlobalOrder(t *testing.T) {
	globalShortcuts := GetByCategory(CategoryGlobal)
	instanceIndex := -1
	serviceIndex := -1
	for i, s := range globalShortcuts {
		if s.Key == "ctrl+n" {
			instanceIndex = i
		}
		if s.Key == "ctrl+s" {
			serviceIndex = i
		}
	}
	assert.NotEqual(t, -1, instanceIndex, "ctrl+n shortcut should be present in global shortcuts")
	assert.NotEqual(t, -1, serviceIndex, "ctrl+s shortcut should be present in global shortcuts")
	assert.Less(t, instanceIndex, serviceIndex, "Instances (ctrl+n) should appear before Services (ctrl+s)")
}

func TestGetByCategory(t *testing.T) {
	// Test Service List Category
	shortcuts := GetByCategory(CategoryServiceList)
	assert.NotEmpty(t, shortcuts)
	for _, s := range shortcuts {
		assert.Equal(t, CategoryServiceList, s.Category)
	}

	// Test Empty Category
	empty := GetByCategory("NonExistent")
	assert.Empty(t, empty)

	// Test Instance List Category
	instShortcuts := GetByCategory(CategoryInstanceList)
	assert.NotEmpty(t, instShortcuts)
	keys := make(map[string]string)
	for _, s := range instShortcuts {
		keys[s.Key] = s.Description
	}
	assert.Equal(t, "Auth", keys["a"])
	assert.Equal(t, "Open URL", keys["o"])
	assert.Equal(t, "Proxy", keys["p"])
}

func TestFormatByCategory(t *testing.T) {
	// Test basic formatting
	// We assume Service List has "r" for Refresh
	s := FormatByCategory(CategoryServiceList, nil)
	assert.Contains(t, s, "[white:black] r [-][black:darkcyan] Refresh [-]")

	// Test formatting with overrides
	// Assume we override "r" to "Reload" (just as an example, though typically we override description based on key)
	// The overrides map maps Key (e.g. "r") to Description (e.g. "Reload")
	overrides := map[string]string{
		"r": "Reload",
	}
	s2 := FormatByCategory(CategoryServiceList, overrides)
	assert.Contains(t, s2, "[white:black] r [-][black:darkcyan] Reload [-]")
	assert.NotContains(t, s2, "Refresh")
}

func TestFormat(t *testing.T) {
	s := Shortcut{Key: "k", Description: "Desc"}
	formatted := s.Format()
	assert.Equal(t, "[dodgerblue]<k> [white]Desc", formatted)
}
