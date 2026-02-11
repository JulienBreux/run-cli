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

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	assert.NotEmpty(t, Registry)

	// Check for a few expected shortcuts
	foundHelp := false
	for _, s := range Registry {
		if s.Key == "?" && s.Rune == '?' {
			foundHelp = true
			break
		}
	}
	assert.True(t, foundHelp, "Registry should contain the help shortcut '?'")
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
