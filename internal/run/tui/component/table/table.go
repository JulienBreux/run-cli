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

package table

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Table represents a Table.
type Table struct {
	Title string
	Table *tview.Table
}

// New creates a new table.
func New(title string) *Table {
	tview.Borders.Horizontal = '-'
	tview.Borders.Vertical = '|'
	tview.Borders.TopLeft = '+'
	tview.Borders.TopRight = '+'
	tview.Borders.BottomLeft = '+'
	tview.Borders.BottomRight = '+'

	table := tview.NewTable().SetBorders(false).SetSelectable(true, false).SetFixed(1, 1)
	table.SetBorder(true)
	table.SetBorderColor(tcell.ColorDarkCyan)
	table.SetBorderAttributes(tcell.AttrBold)
	table.SetBorderPadding(0, 0, 1, 1)

	table.SetTitle(" " + title + " (0) ")
	table.SetTitleColor(tcell.ColorLightCyan)
	table.SetTitleAlign(tview.AlignCenter)

	return &Table{
		Title: title,
		Table: table,
	}
}

// SetHeaders sets the table headers.
// Deprecated: Use SetHeadersWithExpansions instead.
func (t *Table) SetHeaders(headers []string) {
	for i, h := range headers {
		addTableHeader(t.Table, i, h, 1)
	}
}

// SetHeadersWithExpansions sets the table headers with custom expansion values.
func (t *Table) SetHeadersWithExpansions(headers []string, expansions []int) {
	for i, h := range headers {
		exp := 1
		if i < len(expansions) {
			exp = expansions[i]
		}
		addTableHeader(t.Table, i, h, exp)
	}
}

// addTableHeader adds a table header.
func addTableHeader(t *tview.Table, col int, val string, expansion int) {
	t.SetCell(0, col, tview.NewTableCell(val).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorLightCyan).
		SetSelectable(false).
		SetExpansion(expansion))
}
