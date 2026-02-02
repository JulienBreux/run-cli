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

// TODO: Refactor with service and job
package workerpool

import (
	"fmt"
	"strings"

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"NAME",
		"REGION",
		"LAST UPDATED",
		"SCALING",
		"MODIFIED BY",
		"LABELS"}

	listExpansions = []int{
		2, // NAME
		1, // REGION
		2, // LAST UPDATED
		2, // SCALING
		2, // MODIFIED BY
		3, // LABELS
	}

	listTable *table.Table
	workers   []model_workerpool.WorkerPool
)

const (
	LIST_PAGE_TITLE     = "Worker Pools"
	LIST_PAGE_ID        = "workers-list"
	LIST_PAGE_SHORTCUT  = tcell.KeyCtrlW
	SCALE_MODAL_PAGE_ID = "scale-workerpool"
)

var listWorkerPoolsFunc = api_workerpool.List

// List returns a list of workers.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// Load populates the table with the provided list of worker pools.
func Load(newWorkers []model_workerpool.WorkerPool) {
	workers = newWorkers
	render(workers)
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	listTable.Table.SetTitle(fmt.Sprintf(" %s loading ", LIST_PAGE_TITLE))

	// Clear shortcuts
	if footer.ContextShortcutView != nil {
		footer.ContextShortcutView.Clear()
	}

	app.SetFocus(listTable.Table)

	go func() {
		// Fetch real data
		var err error
		workers, err = listWorkerPoolsFunc(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(workers) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
					Shortcuts() // Ensure shortcuts are updated (cleared)
				}
				onResult(err)
			}()

			if err != nil {
				return
			}

			render(workers)
		})
	}()
}

func render(workers []model_workerpool.WorkerPool) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	for i, w := range workers {
		var labels []string
		for k, v := range w.Labels {
			labels = append(labels, fmt.Sprintf("%s: %s", k, v))
		}

		scaling := "n/a"
		if w.Scaling != nil {
			scaling = fmt.Sprintf("Manual: %d", w.Scaling.ManualInstanceCount)
		}

		row := i + 1 // +1 for header row
		listTable.Table.SetCell(row, 0, tview.NewTableCell(w.DisplayName))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(w.Region))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(humanize.Time(w.UpdateTime)))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(scaling))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(w.LastModifier))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(strings.Join(labels, ", ")))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(workers)))

	Shortcuts() // Refresh shortcuts (handles empty list case)
}

// GetSelectedWorkerPool returns the Name and Region of the selected worker pool.
func GetSelectedWorkerPool() (string, string) {
	row, _ := listTable.Table.GetSelection()
	if row < 1 { // Header row or no selection
		return "", ""
	}
	// 0: Name (DisplayName), 1: Region
	name := listTable.Table.GetCell(row, 0).Text
	region := listTable.Table.GetCell(row, 1).Text
	return name, region
}

// GetSelectedWorkerPoolFull returns the full workerpool object for the selected row.
func GetSelectedWorkerPoolFull() *model_workerpool.WorkerPool {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(workers) == 0 {
		return nil
	}
	return &workers[row-1]
}

func Shortcuts() {
	if footer.ContextShortcutView == nil {
		return
	}
	footer.ContextShortcutView.Clear()

	if len(workers) == 0 {
		return
	}

	shortcuts := `[dodgerblue]<r> [white]Refresh  [dodgerblue]<d> [white]Describe  [dodgerblue]<s> [white]Scale`
	footer.ContextShortcutView.SetText(shortcuts)
}
