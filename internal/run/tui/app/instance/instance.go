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

package instance

import (
	"fmt"
	"strings"

	api_instance "github.com/JulienBreux/run-cli/internal/run/api/instance"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/shortcut"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"",
		"AUTH",
		"NAME",
		"REGION",
		"STATUS",
		"CONTAINERS",
		"LAST UPDATED",
	}

	listExpansions = []int{
		1, // PROXY
		1, // AUTH
		2, // NAME
		1, // REGION
		1, // STATUS
		2, // CONTAINERS
		2, // LAST UPDATED
	}

	listTable *table.Table
	instances []model.Instance
)

const (
	LIST_PAGE_TITLE    = "Instances"
	LIST_PAGE_ID       = "instances-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlI
)

var ListInstancesFunc = api_instance.List

// List returns a list table of instances.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// Load populates the table with the provided list of instances.
func Load(newInstances []model.Instance) {
	instances = newInstances
	render(instances)
}

// ListReload reloads instances for the current project and region.
func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.SetTitle(fmt.Sprintf(" %s loading... ", LIST_PAGE_TITLE))

	if footer.ContextShortcutView != nil {
		footer.ContextShortcutView.Clear()
	}

	app.SetFocus(listTable.Table)

	go func() {
		var err error
		instances, err = ListInstancesFunc(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(instances) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
					Shortcuts()
				}
				if onResult != nil {
					onResult(err)
				}
			}()

			if err != nil {
				return
			}

			render(instances)
		})
	}()
}

func render(insts []model.Instance) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	for i, inst := range insts {
		var containerNames []string
		for _, c := range inst.Containers {
			containerNames = append(containerNames, c.Name)
		}
		containersStr := strings.Join(containerNames, ", ")
		if containersStr == "" {
			containersStr = "-"
		}

		status := inst.Status()
		var formattedStatus string
		switch status {
		case "Ready":
			formattedStatus = fmt.Sprintf("[green]%s", status)
		case "Failed":
			formattedStatus = fmt.Sprintf("[red]%s", status)
		case "Pending", "Reconciling":
			formattedStatus = fmt.Sprintf("[yellow]%s", status)
		default:
			formattedStatus = fmt.Sprintf("[grey]%s", status)
		}

		updatedAt := "-"
		if !inst.UpdateTime.IsZero() {
			updatedAt = humanize.Time(inst.UpdateTime)
		}

		proxyStatus := ""
		if inst.Proxy != nil && inst.Proxy.Enabled {
			proxyStatus = "[green]P"
		}

		authStatus := "[red]Yes"
		if inst.InvokerIamDisabled {
			authStatus = "[green]No"
		}

		row := i + 1
		listTable.Table.SetCell(row, 0, tview.NewTableCell(proxyStatus))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(authStatus))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(inst.ID()))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(inst.Region))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(formattedStatus))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(containersStr))
		listTable.Table.SetCell(row, 6, tview.NewTableCell(updatedAt))
	}

	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(insts)))
	Shortcuts()
}

// GetSelectedInstance returns the Name and Region of the selected instance.
func GetSelectedInstance() (string, string) {
	row, _ := listTable.Table.GetSelection()
	if row < 1 {
		return "", ""
	}
	name := listTable.Table.GetCell(row, 2).Text
	region := listTable.Table.GetCell(row, 3).Text
	return name, region
}

// GetSelectedInstanceFull returns the full instance object for the selected row.
func GetSelectedInstanceFull() *model.Instance {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(instances) == 0 {
		return nil
	}
	return &instances[row-1]
}

// Shortcuts updates footer shortcuts for instances list.
func Shortcuts() {
	if footer.ContextShortcutView == nil {
		return
	}
	footer.ContextShortcutView.Clear()

	if len(instances) == 0 {
		return
	}

	s := shortcut.FormatByCategory(shortcut.CategoryInstanceList, nil)
	footer.ContextShortcutView.SetText(s)
}
