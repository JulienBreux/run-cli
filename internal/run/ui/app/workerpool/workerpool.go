// TODO: Refactor with service and job
package workerpool

import (
	"fmt"
	"strings"

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/table"
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
	LIST_PAGE_TITLE    = "Worker Pools"
	LIST_PAGE_ID       = "workers-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlW
)

// List returns a list of workers.
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
		workers, err = api_workerpool.List(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(workers) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
				}
				onResult(err)
			}()

			if err != nil {
				return
			}

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
		})
	}()
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
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<r> [white]Refresh
[dodgerblue]<d> [white]Describe`
	header.ContextShortcutView.SetText(shortcuts)
}
