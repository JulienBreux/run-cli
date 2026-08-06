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

package header

import (
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/shortcut"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
	"github.com/JulienBreux/run-cli/pkg/version"
	"github.com/rivo/tview"
)

var (
	infoView *tview.TextView
)

// New returns a TView header.
// Header composition:
// | Info | Global Shortcuts | Logo |
func New(currentInfo info.Info) *tview.Flex {
	return tview.NewFlex().
		AddItem(columnInfo(currentInfo), 50, 1, false).
		AddItem(columnShortcuts(), 0, 1, false).
		AddItem(logo.New(), 50, 1, false)
}

// UpdateInfo updates the info view.
func UpdateInfo(currentInfo info.Info) {
	infoView.Clear()

	_, _ = fmt.Fprintf(infoView, "[white]Project:        [#bd93f9]%s\n", currentInfo.Project)
	_, _ = fmt.Fprintf(infoView, "[white]Region:         [#bd93f9]%s\n", currentInfo.Region)
	_, _ = fmt.Fprintf(infoView, "[white]User:           [#bd93f9]%s\n", currentInfo.User)
	_, _ = fmt.Fprintf(infoView, "[white]Version:        [#bd93f9]%s\n", version.Version)
}

// returns the info column.
func columnInfo(currentInfo info.Info) *tview.TextView {
	infoView = tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	UpdateInfo(currentInfo)
	return infoView
}

// returns the shortcuts column.
func columnShortcuts() *tview.Flex {
	col1 := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	col2 := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)

	shortcuts := shortcut.GetByCategory(shortcut.CategoryGlobal)
	// Filter out "Esc"
	var visibleShortcuts []shortcut.Shortcut
	for _, s := range shortcuts {
		if s.Key != "Esc" {
			visibleShortcuts = append(visibleShortcuts, s)
		}
	}

	// Split roughly in half or specific logic
	// Registry order: S, J, W, D, ?, P, R, Z, L
	// Desired Col 1: P, R, Z, L (Indices 5, 6, 7, 8)
	// Desired Col 2: S, J, W, D, ? (Indices 0, 1, 2, 3, 4)

	// Since we know the order, let's just pick them.
	// But that defeats the purpose of dynamic registry somewhat.
	// Let's just list them. Using the registry order might change the visual layout.
	// Registry order is: Services, Jobs, Workers, Domain, Help, Project, Region, Console, Releases.

	// If I iterate:
	// Col 1: Services, Jobs, Workers, Domain
	// Col 2: Help, Project, Region, Console, Releases

	// Let's stick to the registry order for now, splitting in half.
	mid := len(visibleShortcuts) / 2
	// Adjust mid to shift more to Col 2 if odd? 9 items. 4 in Col 1, 5 in Col 2.
	// Services, Jobs, Workers, Domain -> Col 1
	// Help, Project, Region, Console, Releases -> Col 2

	// Wait, the current layout has Project/Region in Col 1.
	// I should probably reorder the Registry to match the desired visual layout if I want consistency, OR accept the new order.
	// The prompt implies "using the shortcut package", so I should rely on it.
	// I'll split them evenly.

	for i, s := range visibleShortcuts {
		formatted := s.Format() + "\n"
		if i < mid {
			// This puts Services... in Col 1.
			// Current Col 1 has Project...
			// I'll swap the columns in the Flex addition if I want Services in Col 2?
			// No, standard reading is Col 1 then Col 2.
			// Current: Col 1 (Project) | Col 2 (Services)
			// Registry: Services, ..., Project

			// If I want Project in Col 1, I should reorder Registry or filter specifically.
			// Reordering Registry seems cleanest for "Single Source of Truth defining order".

			_, _ = fmt.Fprint(col2, formatted) // Put early items in Col 2 to match existing (Services on right)
		} else {
			_, _ = fmt.Fprint(col1, formatted) // Put later items in Col 1 (Project on left)
		}
	}

	return tview.NewFlex().
		AddItem(col1, 20, 1, false).
		AddItem(col2, 0, 1, false)
}
