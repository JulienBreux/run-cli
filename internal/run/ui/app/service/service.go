package service

import (
	"fmt"
	"os/exec"
	"runtime"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"SERVICE",
		"REGION",
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

func ListReload(currentInfo info.Info) {
	listTable.Table.Clear()
	listTable.SetHeaders(listHeaders)

	// Fetch real data
	services, err := api_service.List(currentInfo.Project, currentInfo.Region)
	if err != nil {
		// Keep empty if error
		listTable.Table.SetTitle(fmt.Sprintf(" %s (Error) ", LIST_PAGE_TITLE))
		return
	}

	for i, s := range services {
		row := i + 1 // +1 for header row
		listTable.Table.SetCell(row, 0, tview.NewTableCell(s.Name))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(s.Region))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(s.URI))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(s.LastModifier))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(s.UpdateTime.Format("2006-01-02 15:04:05")))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(services)))
}

// GetSelectedServiceURL returns the URL of the currently selected service.
func GetSelectedServiceURL() string {
	row, _ := listTable.Table.GetSelection()
	if row == 0 { // Header row
		return ""
	}
	// URL is now at index 2 (0: Service, 1: Region, 2: URL)
	cell := listTable.Table.GetCell(row, 2)
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
