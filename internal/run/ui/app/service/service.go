package service

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	model "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"SERVICE",
		"URL",
		"LAST DEPLOYED BY",
		"LAST DEPLOYED AT"}

	listTable *table.Table
)

const (
	LIST_PAGE_TITLE    = "Services"
	LIST_PAGE_ID       = "services-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlS
)

// List returns a list of services.
func List() *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeaders(listHeaders)
	return listTable
}

func ListReload() {
	listTable.Table.Clear()
	listTable.SetHeaders(listHeaders)

	// Mock data.
	for i := 0; i < 25; i++ {
		s := model.Service{
			Name:         fmt.Sprintf("service-%02d", i+1),
			URI:          fmt.Sprintf("https://service-%02d-abcdefgh-uc.a.run.app", i+1),
			LastModifier: fmt.Sprintf("user%02d@example.com", i+1),
			UpdateTime:   time.Now().Add(-time.Duration(i) * time.Hour),
		}
		row := i + 1 // +1 for header row
		listTable.Table.SetCell(row, 0, tview.NewTableCell(s.Name))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(s.URI))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(s.LastModifier))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(s.UpdateTime.Format("2006-01-02 15:04:05")))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", listTable.Title, 25))
}

// GetSelectedServiceURL returns the URL of the currently selected service.
func GetSelectedServiceURL() string {
	row, _ := listTable.Table.GetSelection()
	if row == 0 { // Header row
		return ""
	}
	cell := listTable.Table.GetCell(row, 1) // URL is at index 1
	return cell.Text
}

// HandleShortcuts handles service-specific shortcuts.
func HandleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	// Open URL
	if event.Rune() == 'o' {
		url := GetSelectedServiceURL()
		if url != "" {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "linux":
				cmd = exec.Command("xdg-open", url)
			case "windows":
				cmd = exec.Command("cmd", "/c", "start", url)
			case "darwin":
				cmd = exec.Command("open", url)
			default:
				return event // Do nothing if OS is not supported
			}
			_ = cmd.Run() // Ignore error for now, ideally log it
		}
		return nil // Consume the event
	}

	return event
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs
[dodgerblue]<s> [white]Scale
[dodgerblue]<o> [white]Open URL`
	header.ContextShortcutView.SetText(shortcuts)
}
