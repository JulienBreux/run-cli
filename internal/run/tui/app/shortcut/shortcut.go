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

import "github.com/gdamore/tcell/v2"

// Shortcut represents a keyboard shortcut.
type Shortcut struct {
	Key         string
	TCellKey    tcell.Key
	Rune        rune
	Description string
	Category    string
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
	{Key: "Enter", TCellKey: tcell.KeyEnter, Description: "Service Dashboard", Category: CategoryServiceList},
	{Key: "r", Rune: 'r', Description: "Refresh List", Category: CategoryServiceList},
	{Key: "l", Rune: 'l', Description: "View Logs", Category: CategoryServiceList},
	{Key: "d", Rune: 'd', Description: "Describe Service", Category: CategoryServiceList},
	{Key: "s", Rune: 's', Description: "Scale Service", Category: CategoryServiceList},
	{Key: "a", Rune: 'a', Description: "Manage Authentication", Category: CategoryServiceList},
	{Key: "t", Rune: 't', Description: "Traffic Split", Category: CategoryServiceList},
	{Key: "o", Rune: 'o', Description: "Open URL", Category: CategoryServiceList},
	{Key: "p", Rune: 'p', Description: "Toggle Proxy", Category: CategoryServiceList},

	// Service Dashboard
	{Key: "Tab", TCellKey: tcell.KeyTab, Description: "Next Tab", Category: CategoryServiceDashboard},
	{Key: "Shift+Tab", TCellKey: tcell.KeyBacktab, Description: "Previous Tab", Category: CategoryServiceDashboard},
	{Key: "t", Rune: 't', Description: "Traffic Split", Category: CategoryServiceDashboard},

	// Job List
	{Key: "Enter", TCellKey: tcell.KeyEnter, Description: "Job Dashboard", Category: CategoryJobList},
	{Key: "r", Rune: 'r', Description: "Refresh List", Category: CategoryJobList},
	{Key: "l", Rune: 'l', Description: "View Logs", Category: CategoryJobList},
	{Key: "d", Rune: 'd', Description: "Describe Job", Category: CategoryJobList},
	{Key: "x", Rune: 'x', Description: "Execute Job", Category: CategoryJobList},

	// Worker List
	{Key: "r", Rune: 'r', Description: "Refresh List", Category: CategoryWorkerList},
	{Key: "d", Rune: 'd', Description: "Describe Worker Pool", Category: CategoryWorkerList},
	{Key: "s", Rune: 's', Description: "Scale Worker Pool", Category: CategoryWorkerList},

	// Domain Mapping
	{Key: "Enter", TCellKey: tcell.KeyEnter, Description: "DNS Info", Category: CategoryDomainMapping},
	{Key: "r", Rune: 'r', Description: "Refresh List", Category: CategoryDomainMapping},
	{Key: "o", Rune: 'o', Description: "Open URL", Category: CategoryDomainMapping},
}
