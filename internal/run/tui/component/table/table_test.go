package table

import (
	"testing"
)

func TestNew(t *testing.T) {
	title := "Test Table"
	table := New(title)

	if table.Title != title {
		t.Errorf("Expected title '%s', got '%s'", title, table.Title)
	}
	if table.Table == nil {
		t.Error("Expected table to be initialized, got nil")
	}
}

func TestSetHeadersWithExpansions(t *testing.T) {
	title := "Test Table"
	tbl := New(title)
	headers := []string{"Col1", "Col2"}
	expansions := []int{1, 2}

	tbl.SetHeadersWithExpansions(headers, expansions)

	// Since we cannot easily inspect the internal state of tview.Table via public API (cells),
	// we rely on the fact that no panic occurred and basic property checks.
	// tview.Table doesn't expose a way to get cell content easily without SetCell.
	// Wait, GetCell exists.
	
	cell := tbl.Table.GetCell(0, 0)
	if cell.Text != "Col1" {
		t.Errorf("Expected header 1 'Col1', got '%s'", cell.Text)
	}

	cell2 := tbl.Table.GetCell(0, 1)
	if cell2.Text != "Col2" {
		t.Errorf("Expected header 2 'Col2', got '%s'", cell2.Text)
	}
}

func TestSetHeaders(t *testing.T) {
	title := "Test Table"
	tbl := New(title)
	headers := []string{"A", "B"}

	tbl.SetHeaders(headers)

	cell := tbl.Table.GetCell(0, 0)
	if cell.Text != "A" {
		t.Errorf("Expected header 'A', got '%s'", cell.Text)
	}
}
