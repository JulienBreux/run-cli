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

// TODO: Refactor with service and worker pool
package job

import (
	"fmt"
	"strings"

	api_job "github.com/JulienBreux/run-cli/internal/run/api/job"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
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

var listJobsFunc = api_job.List

// List returns a list of jobs.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	app.SetFocus(listTable.Table)
	return listTable
}

// Load populates the table with the provided list of jobs.
func Load(newJobs []model_job.Job) {
	jobs = newJobs
	render(jobs)
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	listTable.Table.SetTitle(fmt.Sprintf(" %s loading ", LIST_PAGE_TITLE))

	app.SetFocus(listTable.Table)

	go func() {
		// Fetch real data
		var err error
		jobs, err = listJobsFunc(currentInfo.Project, currentInfo.Region)

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

			render(jobs)
		})
	}()
}

func render(jobs []model_job.Job) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

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
	footer.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<r> [white]Refresh  [dodgerblue]<d> [white]Describe  [dodgerblue]<l> [white]Logs  [dodgerblue]<x> [white]Execute  [dodgerblue]<enter> [white]Details`
	footer.ContextShortcutView.SetText(shortcuts)
}
