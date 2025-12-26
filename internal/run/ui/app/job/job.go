// TODO: Refactor with service and worker pool
package job

import (
	"fmt"
	"strings"

	api_job "github.com/JulienBreux/run-cli/internal/run/api/job"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/dustin/go-humanize"
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

	listExpansions = []int{
		2, // NAME
		2, // STATUS OF LAST EXECUTION
		2, // LAST EXECUTED
		1, // REGION
		2, // CREATED BY
	}

	listTable *table.Table
	jobs      []model_job.Job
)

const (
	LIST_PAGE_TITLE    = "Jobs"
	LIST_PAGE_ID       = "jobs-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlJ
)

// List returns a list of jobs.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	app.SetFocus(listTable.Table)
	return listTable
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	app.SetFocus(listTable.Table)

	go func() {
		// Fetch real data
		var err error
		jobs, err = api_job.List(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(jobs) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
				}
				onResult(err)
			}()

			if err != nil {
				// listTable.Table.SetTitle(fmt.Sprintf(" %s (Error) ", LIST_PAGE_TITLE)) // Removed error from title
				return
			}

			for i, j := range jobs {
				// Extract info
				nameParts := strings.Split(j.Name, "/")
				displayName := nameParts[len(nameParts)-1]

				status := "-"
				if j.TerminalCondition != nil {
					status = j.TerminalCondition.State
				}

				lastExecuted := "-"
				if j.LatestCreatedExecution != nil {
					lastExecuted = humanize.Time(j.LatestCreatedExecution.CreateTime)
				}

				row := i + 1 // +1 for header row
				listTable.Table.SetCell(row, 0, tview.NewTableCell(displayName))
				listTable.Table.SetCell(row, 1, tview.NewTableCell(status))
				listTable.Table.SetCell(row, 2, tview.NewTableCell(lastExecuted))
				listTable.Table.SetCell(row, 3, tview.NewTableCell(j.Region))
				listTable.Table.SetCell(row, 4, tview.NewTableCell(j.Creator))
			}

			// Refresh title
			listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(jobs)))
		})
	}()
}

// GetSelectedJob returns the Name and Region of the selected job.
func GetSelectedJob() (string, string) {
	row, _ := listTable.Table.GetSelection()
	if row < 1 { // Header row or no selection
		return "", ""
	}
	// 0: Name, 3: Region
	name := listTable.Table.GetCell(row, 0).Text
	region := listTable.Table.GetCell(row, 3).Text
	return name, region
}

// GetSelectedJobFull returns the full job object for the selected row.
func GetSelectedJobFull() *model_job.Job {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(jobs) == 0 {
		return nil
	}
	return &jobs[row-1]
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs`
	header.ContextShortcutView.SetText(shortcuts)
}
