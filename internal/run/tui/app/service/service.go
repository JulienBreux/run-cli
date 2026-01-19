// TODO: Refactor with job and worker pool
package service

import (
	"fmt"
	"os"
	"strings"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/pkg/browser"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"SERVICE",
		"REGION",
		"SCALING",
		"URL",
		"LAST DEPLOYED BY",
		"LAST DEPLOYED AT"}

	listExpansions = []int{
		2, // SERVICE
		1, // REGION
		1, // SCALING
		4, // URL
		2, // LAST DEPLOYED BY
		1, // LAST DEPLOYED AT
	}

	listTable *table.Table
	services  []model_service.Service
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
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// Load populates the table with the provided list of services.
func Load(newServices []model_service.Service) {
	services = newServices
	render(services)
}

func ListReload(app *tview.Application, currentInfo info.Info, onResult func(error)) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
	listTable.Table.SetTitle(fmt.Sprintf(" %s loading ", LIST_PAGE_TITLE))

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
				}
				onResult(err)
			}()

			if err != nil {
				// Keep empty if error
				return
			}

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

		listTable.Table.SetCell(row, 0, tview.NewTableCell(s.Name))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(s.Region))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(scaling))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(s.URI))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(s.LastModifier))
		listTable.Table.SetCell(row, 5, tview.NewTableCell(humanize.Time(s.UpdateTime)))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(svc)))
}

// GetSelectedServiceURL returns the URL of the currently selected service.
func GetSelectedServiceURL() string {
	row, _ := listTable.Table.GetSelection()
	if row == 0 { // Header row
		return ""
	}
	// URL is now at index 3 (0: Service, 1: Region, 2: Scaling, 3: URL)
	cell := listTable.Table.GetCell(row, 3)
	return cell.Text
}

// GetSelectedService returns the Name and Region of the selected service.
func GetSelectedService() (string, string) {
	row, _ := listTable.Table.GetSelection()
	if row < 1 { // Header row or no selection
		return "", ""
	}
	// 0: Service, 1: Region
	name := listTable.Table.GetCell(row, 0).Text
	region := listTable.Table.GetCell(row, 1).Text
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
		u := GetSelectedServiceURL()
		if u != "" && !strings.HasSuffix(os.Args[0], ".test") {
			_ = browser.OpenURL(u)
			return event
		}
		return nil // Consume the event
	}

	return event
}

func Shortcuts() {
	footer.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<r> [white]Refresh  [dodgerblue]<d> [white]Describe  [dodgerblue]<l> [white]Logs  [dodgerblue]<s> [white]Scale  [dodgerblue]<a> [white]Auth  [dodgerblue]<o> [white]Open URL  [dodgerblue]<enter> [white]Details`
	footer.ContextShortcutView.SetText(shortcuts)
}
