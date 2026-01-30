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

// TODO: Refactor with job and worker pool
package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/proxy"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"",
		"SERVICE",
		"REGION",
		"SCALING",
		"AUTH",
		"URL",
		"LAST DEPLOYED BY",
		"LAST DEPLOYED AT"}

	listExpansions = []int{
		1, // PROXY
		2, // SERVICE
		1, // REGION
		1, // SCALING
		1, // AUTH
		4, // URL
		2, // LAST DEPLOYED BY
		1, // LAST DEPLOYED AT
	}

	listTable *table.Table
	services  []model_service.Service
	proxies   *proxy.Manager
)

const (
	LIST_PAGE_TITLE     = "Services"
	LIST_PAGE_ID        = "services-list"
	LIST_PAGE_SHORTCUT  = tcell.KeyCtrlS
	SCALE_MODAL_PAGE_ID = "scale"
)

var listServicesFunc = api_service.List

// Fetch retrieves the list of services from the API.
func Fetch(projectID, region string) ([]model_service.Service, error) {
	return listServicesFunc(projectID, region)
}

// List returns a list of services.
func List(app *tview.Application) *table.Table {
	proxies = proxy.NewManager()
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// syncProxies updates the service list with running proxy information.
func syncProxies(svcs []model_service.Service) {
	if proxies == nil {
		return
	}
	for i := range svcs {
		if info := proxies.GetInfo(svcs[i].Name); info != nil {
			svcs[i].Proxy = &model_service.ProxyStatus{
				Enabled: true,
				Port:    info.Port,
				URL:     fmt.Sprintf("http://127.0.0.1:%d", info.Port),
			}
		}
	}
}

// Load populates the table with the provided list of services.
func Load(newServices []model_service.Service) {
	services = newServices
	syncProxies(services)
	render(services)
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
		services, err = Fetch(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(services) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
					Shortcuts() // Ensure shortcuts are updated (cleared)
				}
				onResult(err)
			}()

			if err != nil {
				// Keep empty if error
				return
			}

			syncProxies(services)
			render(services)
		})
	}()
}

func render(svc []model_service.Service) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	for i, s := range svc {
		row := i + 1 // +1 for header row

		scaling := "n/a"
		if s.Scaling != nil {
			switch s.Scaling.ScalingMode {
			case "AUTOMATIC":
				scaling = fmt.Sprintf("Auto: min %d", s.Scaling.MinInstances)
				if s.Scaling.MaxInstances != 0 {
					scaling += fmt.Sprintf(", max %d", s.Scaling.MaxInstances)
				}
			case "MANUAL":
				scaling = fmt.Sprintf("Manual: %d", s.Scaling.ManualInstanceCount)
			}
		}

		proxyStatus := " "
		if s.Proxy != nil && s.Proxy.Enabled {
			proxyStatus = "[green]P[white]"
		}

		authStatus := "[red]Yes"
		if s.Security != nil && s.Security.InvokerIAMDisabled {
			authStatus = "[green]No"
		}

		listTable.Table.SetCell(row, 0, tview.NewTableCell(proxyStatus))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(s.Name))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(s.Region))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(scaling))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(authStatus))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(s.URI))
		listTable.Table.SetCell(row, 6, tview.NewTableCell(s.LastModifier))
		listTable.Table.SetCell(row, 7, tview.NewTableCell(humanize.Time(s.UpdateTime)))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(svc)))

	// selection change
	listTable.Table.SetSelectionChangedFunc(func(row, column int) {
		Shortcuts()
	})

	Shortcuts() // Refresh shortcuts (handles empty list case)
}

// GetSelectedServiceURL returns the URL of the currently selected service.
func GetSelectedServiceURL() string {
	row, _ := listTable.Table.GetSelection()
	if row == 0 { // Header row
		return ""
	}
	// URL is now at index 5 (0: Proxy, 1: Service, 2: Region, 3: Scaling, 4: Auth, 5: URL)
	cell := listTable.Table.GetCell(row, 5)
	return cell.Text
}

// GetSelectedService returns the Name and Region of the selected service.
func GetSelectedService() (string, string) {
	row, _ := listTable.Table.GetSelection()
	if row < 1 { // Header row or no selection
		return "", ""
	}
	// 1: Service, 2: Region
	name := listTable.Table.GetCell(row, 1).Text
	region := listTable.Table.GetCell(row, 2).Text
	return name, region
}

// GetSelectedServiceFull returns the full service object for the selected row.
func GetSelectedServiceFull() *model_service.Service {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(services) == 0 {
		return nil
	}
	return &services[row-1]
}

// HandleShortcuts handles service-specific shortcuts.
func HandleShortcuts(event *tcell.EventKey) *tcell.EventKey {
	// Open URL
	if event.Rune() == 'o' {
		url := GetSelectedServiceURL()
		svc := GetSelectedServiceFull()
		if svc != nil && svc.Proxy != nil && svc.Proxy.Enabled {
			url = svc.Proxy.URL
		}

		if url != "" && !strings.HasSuffix(os.Args[0], ".test") {
			_ = browser.OpenURL(url)
			return event
		}
		return nil // Consume the event
	}

	// Toggle Proxy
	if event.Rune() == 'p' {
		toggleProxy()
		return nil
	}

	return event
}

func toggleProxy() {
	svc := GetSelectedServiceFull()
	if svc == nil {
		return
	}

	if svc.Proxy != nil && svc.Proxy.Enabled {
		// Stop Proxy
		_ = proxies.Stop(svc.Name)
		svc.Proxy.Enabled = false
		svc.Proxy.Port = 0
		svc.Proxy.URL = ""
	} else {
		// Start Proxy
		info, err := proxies.Start(context.Background(), svc.Name, svc.URI)
		if err != nil {
			// Handle error
			return
		}
		svc.Proxy = &model_service.ProxyStatus{
			Enabled: true,
			Port:    info.Port,
			URL:     fmt.Sprintf("http://127.0.0.1:%d", info.Port),
		}
	}
	render(services)
	Shortcuts() // Ensure footer updates immediately
}

func Shortcuts() {
	if footer.ContextShortcutView == nil {
		return
	}
	footer.ContextShortcutView.Clear()

	if len(services) == 0 {
		return
	}

	shortcuts := `[dodgerblue]<r> [white]Refresh  [dodgerblue]<d> [white]Describe  [dodgerblue]<l> [white]Logs [dodgerblue]<s> [white]Scale [dodgerblue]<a> [white]Auth [dodgerblue]<o> [white]Open URL  [dodgerblue]<enter> [white]Details`

	// Check selected service proxy status
	svc := GetSelectedServiceFull()
	if svc != nil && svc.Proxy != nil && svc.Proxy.Enabled {
		s := fmt.Sprintf("[white]Auth [dodgerblue]<p> [green]Proxy (127.0.0.1:%d)", svc.Proxy.Port)
		shortcuts = strings.Replace(shortcuts, "[white]Auth", s, 1)
		shortcuts = strings.Replace(shortcuts, "[white]Open URL", "[white]Open URL (proxy)", 1)
	} else {
		shortcuts = strings.Replace(shortcuts, "[white]Auth", "[white]Auth [dodgerblue]<p> [white]Proxy", 1)
	}

	footer.ContextShortcutView.SetText(shortcuts)
}
