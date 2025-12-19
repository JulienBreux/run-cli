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

// setTableRow sets a table row.
func setTableRow(t *tview.Table, row int, statusColor tcell.Color, values ...string) {
	for i, v := range values {
		color := tcell.ColorWhite
		if i == 2 {
			color = statusColor
		}
		t.SetCell(row, i, tview.NewTableCell(v).SetTextColor(color))
	}
}