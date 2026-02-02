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
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Shortcut represents a keyboard shortcut.
type Shortcut struct {
	Key         string
	TCellKey    tcell.Key
	Rune        rune
	Description string
	Category    string
}

// Format returns the formatted string for the shortcut.
func (s Shortcut) Format() string {
	return fmt.Sprintf("[dodgerblue]<%s> [white]%s", s.Key, s.Description)
}

const (
	CategoryGlobal           = "Global"
	CategoryServiceList      = "Service List"
	CategoryServiceDashboard = "Service Dashboard"
	CategoryJobList          = "Job List"
	CategoryWorkerList       = "Worker List"
	CategoryDomainMapping    = "Domain Mapping"
)

var Registry = []Shortcut{
	// Global
	{Key: "Ctrl+S", TCellKey: tcell.KeyCtrlS, Description: "Services List", Category: CategoryGlobal},
	{Key: "Ctrl+J", TCellKey: tcell.KeyCtrlJ, Description: "Jobs List", Category: CategoryGlobal},
	{Key: "Ctrl+W", TCellKey: tcell.KeyCtrlW, Description: "Worker Pools List", Category: CategoryGlobal},
	{Key: "Ctrl+D", TCellKey: tcell.KeyCtrlD, Description: "Domain Mappings List", Category: CategoryGlobal},
	{Key: "?", Rune: '?', Description: "Help Page", Category: CategoryGlobal},
	{Key: "Ctrl+P", TCellKey: tcell.KeyCtrlP, Description: "Project Selection", Category: CategoryGlobal},
	{Key: "Ctrl+R", TCellKey: tcell.KeyCtrlR, Description: "Region Selection", Category: CategoryGlobal},
	{Key: "Esc", TCellKey: tcell.KeyEscape, Description: "Back / Close Modal", Category: CategoryGlobal},

	// Service List
	{Key: "r", Rune: 'r', Description: "Refresh List", Category: CategoryServiceList},
	{Key: "d", Rune: 'd', Description: "Describe Service", Category: CategoryServiceList},
	{Key: "l", Rune: 'l', Description: "Logs", Category: CategoryServiceList},
	{Key: "s", Rune: 's', Description: "Scale", Category: CategoryServiceList},
	{Key: "a", Rune: 'a', Description: "Auth", Category: CategoryServiceList},
	{Key: "t", Rune: 't', Description: "Traffic", Category: CategoryServiceList},
	{Key: "o", Rune: 'o', Description: "Open URL", Category: CategoryServiceList},
	{Key: "enter", TCellKey: tcell.KeyEnter, Description: "Details", Category: CategoryServiceList},
	{Key: "?", Rune: '?', Description: "Help", Category: CategoryServiceList},

	// Service Dashboard
	{Key: "esc", TCellKey: tcell.KeyEscape, Description: "Back", Category: CategoryServiceDashboard},
	{Key: "tab", TCellKey: tcell.KeyTab, Description: "Next Tab", Category: CategoryServiceDashboard},
	{Key: "shift-tab", TCellKey: tcell.KeyBacktab, Description: "Prev Tab", Category: CategoryServiceDashboard},
	{Key: "t", Rune: 't', Description: "Traffic Split", Category: CategoryServiceDashboard},
	{Key: "?", Rune: '?', Description: "Help", Category: CategoryServiceDashboard},

	// Job List
	{Key: "r", Rune: 'r', Description: "Refresh", Category: CategoryJobList},
	{Key: "d", Rune: 'd', Description: "Describe", Category: CategoryJobList},
	{Key: "l", Rune: 'l', Description: "Logs", Category: CategoryJobList},
	{Key: "x", Rune: 'x', Description: "Execute", Category: CategoryJobList},
	{Key: "enter", TCellKey: tcell.KeyEnter, Description: "Details", Category: CategoryJobList},
	{Key: "?", Rune: '?', Description: "Help", Category: CategoryJobList},

	// Worker List
	{Key: "r", Rune: 'r', Description: "Refresh", Category: CategoryWorkerList},
	{Key: "d", Rune: 'd', Description: "Describe", Category: CategoryWorkerList},
	{Key: "s", Rune: 's', Description: "Scale", Category: CategoryWorkerList},
	{Key: "?", Rune: '?', Description: "Help", Category: CategoryWorkerList},

	// Domain Mapping
	{Key: "r", Rune: 'r', Description: "Refresh", Category: CategoryDomainMapping},
	{Key: "o", Rune: 'o', Description: "Open URL", Category: CategoryDomainMapping},
	{Key: "enter", TCellKey: tcell.KeyEnter, Description: "Info", Category: CategoryDomainMapping},
	{Key: "?", Rune: '?', Description: "Help", Category: CategoryDomainMapping},
}

// GetByCategory returns shortcuts for a given category.
func GetByCategory(category string) []Shortcut {
	var filtered []Shortcut
	for _, s := range Registry {
		if s.Category == category {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// FormatByCategory returns a formatted string of shortcuts for a given category.
// overrides is a map of Key -> Description to override default descriptions.
func FormatByCategory(category string, overrides map[string]string) string {
	shortcuts := GetByCategory(category)
	var formatted []string
	for _, s := range shortcuts {
		desc := s.Description
		if val, ok := overrides[s.Key]; ok {
			desc = val
		}
		// Use local struct to format with override
		temp := Shortcut{Key: s.Key, Description: desc}
		formatted = append(formatted, temp.Format())
	}
	return strings.Join(formatted, "  ")
}