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
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/tui/app/shortcut"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "modal-help"
)

// Help represents the help modal component.
type Help struct {
	*tview.Grid
	Table *tview.Table
}

// HelpModal creates a modal to display all keyboard shortcuts.
func HelpModal(app *tview.Application, closeFunc func()) *Help {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false)

	table.SetBorder(true).SetTitle(" Keyboard Shortcuts ")

	row := 0
	currentCategory := ""
	for _, s := range shortcut.Registry {
		if s.Category != currentCategory {
			if row > 0 {
				row++
			}
			table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("[yellow::b]%s", s.Category)).SetSelectable(false))
			currentCategory = s.Category
			row++
		}

		table.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("[dodgerblue]%s", s.Key)).SetSelectable(false))
		table.SetCell(row, 1, tview.NewTableCell(s.Description).SetSelectable(false))
		row++
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyEnter || event.Rune() == '?' {
			closeFunc()
			return nil
		}
		return event
	})

	grid := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 35, 0).
		AddItem(table, 1, 1, 1, 1, 0, 0, true)

	return &Help{
		Grid:  grid,
		Table: table,
	}
}
