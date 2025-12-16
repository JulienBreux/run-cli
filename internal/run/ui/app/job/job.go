package job

import (
	"fmt"
	"strings"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	model "github.com/JulienBreux/run-cli/internal/run/model/job"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"NAME",
		"STATUS OF LAST EXECUTION",
		"LAST EXECUTED",
		"REGION",
		"CREATED BY"}

	listTable *table.Table
)

const (
	LIST_PAGE_TITLE    = "Jobs"
	LIST_PAGE_ID       = "jobs-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlJ
)

// List returns a list of jobs.
func List() *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeaders(listHeaders)
	return listTable
}

func ListReload() {
	listTable.Table.Clear()
	listTable.SetHeaders(listHeaders)

	// Mock data.
	for i := 0; i < 5; i++ {
		region := "us-central1"
		simpleName := fmt.Sprintf("job-processor-%02d", i+1)
		fullName := fmt.Sprintf("projects/test-project/locations/%s/jobs/%s", region, simpleName)

		j := model.Job{
			Name: fullName,
			LatestCreatedExecution: &model.ExecutionReference{
				CreateTime: time.Now().Add(-time.Duration(i*2) * time.Hour),
			},
			TerminalCondition: &condition.Condition{
				State: "Succeeded",
			},
			Creator: fmt.Sprintf("dev-%02d@example.com", i+1),
		}

		if i%3 == 0 {
			j.TerminalCondition.State = "Failed"
		}

		// Extract info
		nameParts := strings.Split(j.Name, "/")
		displayName := nameParts[len(nameParts)-1]
		displayRegion := nameParts[3]

		row := i + 1 // +1 for header row
		listTable.Table.SetCell(row, 0, tview.NewTableCell(displayName))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(j.TerminalCondition.State))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(j.LatestCreatedExecution.CreateTime.Format("2006-01-02 15:04:05")))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(displayRegion))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(j.Creator))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", listTable.Title, 5))
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs
[dodgerblue]<x> [white]Execute`
	header.ContextShortcutView.SetText(shortcuts)
}
