package worker

import (
	"fmt"
	"strings"

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"NAME",
		"STATUS",
		"REGION",
		"LAST UPDATE",
		"MACHINE TYPE",
		"LABELS"}

	listExpansions = []int{
		2, // NAME
		1, // STATUS
		1, // REGION
		2, // LAST UPDATE
		2, // MACHINE TYPE
		3, // LABELS
	}

	listTable *table.Table
)

const (
	LIST_PAGE_TITLE    = "Worker Pools"
	LIST_PAGE_ID       = "workers-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlW
)

// List returns a list of workers.
func List() *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	return listTable
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	go func() {
		// Fetch real data
		workers, err := api_workerpool.List(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer onResult(err)

			if err != nil {
				return
			}

			for i, w := range workers {
				machineType := "-"
				if w.WorkerConfig != nil {
					machineType = w.WorkerConfig.MachineType
				}

				var labels []string
				for k, v := range w.Labels {
					labels = append(labels, fmt.Sprintf("%s=%s", k, v))
				}

				row := i + 1 // +1 for header row
				listTable.Table.SetCell(row, 0, tview.NewTableCell(w.DisplayName))
				listTable.Table.SetCell(row, 1, tview.NewTableCell(w.State))
				listTable.Table.SetCell(row, 2, tview.NewTableCell(w.Region))
				listTable.Table.SetCell(row, 3, tview.NewTableCell(humanize.Time(w.UpdateTime)))
				listTable.Table.SetCell(row, 4, tview.NewTableCell(machineType))
				listTable.Table.SetCell(row, 5, tview.NewTableCell(strings.Join(labels, ", ")))
			}

			// Refresh title
			listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(workers)))
		})
	}()
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs
[dodgerblue]<s> [white]Scale`
	header.ContextShortcutView.SetText(shortcuts)
}
