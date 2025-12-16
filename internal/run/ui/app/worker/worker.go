package worker

import (
	"fmt"
	"strings"

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/table"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"NAME",
		"DEPLOYMENT TYPE",
		"REGION",
		"LAST UPDATE",
		"SCALING",
		"LABELS"}

	listTable *table.Table
)

const (
	LIST_PAGE_TITLE    = "Worker Pools"
	LIST_PAGE_ID       = "worker-pools-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlW
)

// List returns a list of worker pools.
func List() *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeaders(listHeaders)
	return listTable
}

func ListReload(currentInfo info.Info) {
	listTable.Table.Clear()
	listTable.SetHeaders(listHeaders)

	// Fetch real data
	workerPools, err := api_workerpool.List(currentInfo.Project, currentInfo.Region)
	if err != nil {
		listTable.Table.SetTitle(fmt.Sprintf(" %s (Error: %s) ", LIST_PAGE_TITLE, err))
		return
	}

	for i, w := range workerPools {
		// Infer Deployment Type
		deploymentType := "Standard"
		if w.PrivatePoolVpcConfig != nil {
			deploymentType = "Private"
		} else if w.NetworkConfig != nil && w.NetworkConfig.EgressOption == "PRIVATE_ENDPOINT" {
			deploymentType = "Hybrid"
		}

		// Infer Scaling
		scaling := "-"
		if w.WorkerConfig != nil {
			scaling = fmt.Sprintf("Machine: %s, Disk: %dGB", w.WorkerConfig.MachineType, w.WorkerConfig.DiskSizeGb)
		}

		// Format labels
		var labels []string
		for k, v := range w.Labels {
			labels = append(labels, fmt.Sprintf("%s:%s", k, v))
		}

		row := i + 1 // +1 for header row
		listTable.Table.SetCell(row, 0, tview.NewTableCell(w.DisplayName))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(deploymentType))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(w.Region))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(w.UpdateTime.Format("2006-01-02 15:04:05")))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(scaling))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(strings.Join(labels, ", ")))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(workerPools)))
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs
[dodgerblue]<s> [white]Scale`
	header.ContextShortcutView.SetText(shortcuts)
}
