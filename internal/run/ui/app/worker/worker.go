package worker

import (
	"fmt"
	"strings"
	"time"

	model "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
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

func ListReload() {
	listTable.Table.Clear()
	listTable.SetHeaders(listHeaders)

	// Mock data.
	for i := 0; i < 20; i++ {
		region := "us-central1"
		workerPoolName := fmt.Sprintf("my-workerpool-%02d", i+1)
		fullName := fmt.Sprintf("projects/test-project/locations/%s/workerPools/%s", region, workerPoolName)

		w := model.WorkerPool{
			Name:        fullName,
			DisplayName: workerPoolName,
			Region:      region,
			UpdateTime:  time.Now().Add(-time.Duration(i*5) * time.Hour),
			WorkerConfig: &model.WorkerConfig{
				MachineType: "e2-medium",
				DiskSizeGb:  300,
			},
			NetworkConfig: &model.NetworkConfig{
				EgressOption: "PRIVATE_ENDPOINT",
			},
			Labels: map[string]string{"env": "prod", "team": "backend"},
		}

		if i%2 == 0 {
			w.WorkerConfig.MachineType = "e2-small"
			w.NetworkConfig.EgressOption = "NO_EXTERNAL_IP"
			w.Labels["env"] = "dev"
		}

		// Infer Deployment Type
		deploymentType := "Container"

		// Infer Scaling
		scaling := fmt.Sprintf("Manual: 1")

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
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", listTable.Title, 20))
}

func Shortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<d> [white]Describe
[dodgerblue]<l> [white]Logs
[dodgerblue]<s> [white]Scale`
	header.ContextShortcutView.SetText(shortcuts)
}
