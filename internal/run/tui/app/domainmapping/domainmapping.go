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

package domainmapping

import (
	"fmt"
	"strings"

	api_domainmapping "github.com/JulienBreux/run-cli/internal/run/api/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_domainmapping "github.com/JulienBreux/run-cli/internal/run/model/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/shortcut"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	listHeaders = []string{
		"DOMAIN",
		"MAPPED TO",
		"REGION",
		"ADDED BY",
		"CREATED"}

	listExpansions = []int{
		2, // DOMAIN
		2, // MAPPED TO
		1, // REGION
		2, // ADDED BY
		2, // CREATED
	}

	listTable      *table.Table
	domainMappings []model_domainmapping.DomainMapping
)

const (
	LIST_PAGE_TITLE    = "Domain Mappings"
	LIST_PAGE_ID       = "domainmappings-list"
	LIST_PAGE_SHORTCUT = tcell.KeyCtrlD
	MODAL_PAGE_ID      = "modal-dns-records"
)

var listDomainMappingsFunc = api_domainmapping.List

// List returns a list of domain mappings.
func List(app *tview.Application) *table.Table {
	listTable = table.New(LIST_PAGE_TITLE)
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	app.SetFocus(listTable.Table)

	return listTable
}

// Load populates the table with the provided list of domain mappings.
func Load(newDomainMappings []model_domainmapping.DomainMapping) {
	domainMappings = newDomainMappings
	render(domainMappings)
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
		domainMappings, err = listDomainMappingsFunc(currentInfo.Project, currentInfo.Region)

		app.QueueUpdateDraw(func() {
			defer func() {
				if len(domainMappings) == 0 {
					listTable.Table.Clear()
					listTable.SetHeadersWithExpansions(listHeaders, listExpansions)
					Shortcuts() // Ensure shortcuts are updated (cleared)
				}
				onResult(err)
			}()

			if err != nil {
				return
			}

			render(domainMappings)
		})
	}()
}

func render(dms []model_domainmapping.DomainMapping) {
	listTable.Table.Clear()
	listTable.SetHeadersWithExpansions(listHeaders, listExpansions)

	for i, dm := range dms {
		row := i + 1 // +1 for header row

		listTable.Table.SetCell(row, 0, tview.NewTableCell(dm.Name))
		listTable.Table.SetCell(row, 1, tview.NewTableCell(dm.RouteName))
		listTable.Table.SetCell(row, 2, tview.NewTableCell(dm.Region))
		listTable.Table.SetCell(row, 3, tview.NewTableCell(dm.Creator))
		listTable.Table.SetCell(row, 4, tview.NewTableCell(humanize.Time(dm.CreateTime)))
	}

	// Refresh title
	listTable.Table.SetTitle(fmt.Sprintf(" %s (%d) ", LIST_PAGE_TITLE, len(dms)))

	Shortcuts() // Refresh shortcuts (handles empty list case)
}

// GetSelectedDomainMappingFull returns the full domain mapping object for the selected row.
func GetSelectedDomainMappingFull() *model_domainmapping.DomainMapping {
	row, _ := listTable.Table.GetSelection()
	if row < 1 || len(domainMappings) == 0 {
		return nil
	}
	return &domainMappings[row-1]
}

// GetSelectedDomainURL returns the URL of the currently selected domain mapping.
func GetSelectedDomainURL() string {
	row, _ := listTable.Table.GetSelection()
	if row < 1 { // Header row or no selection
		return ""
	}
	// 0: Domain
	domain := listTable.Table.GetCell(row, 0).Text
	return fmt.Sprintf("https://%s", domain)
}

func Shortcuts() {
	if footer.ContextShortcutView == nil {
		return
	}
	footer.ContextShortcutView.Clear()

	if len(domainMappings) == 0 {
		return
	}

	s := shortcut.FormatByCategory(shortcut.CategoryDomainMapping, nil)
	footer.ContextShortcutView.SetText(s)
}

// DomainMappingInfoModal creates a modal to display info for the domain mapping.
func DomainMappingInfoModal(app *tview.Application, dm *model_domainmapping.DomainMapping, closeFunc func()) tview.Primitive {
	var sb strings.Builder

	// Status Section
	fmt.Fprintln(&sb, "[yellow::b]Status[white::-]")
	state := "Unknown"
	message := ""
	for _, c := range dm.Conditions {
		if c.Type == "Ready" {
			if c.State == "True" {
				state = "Ready"
			} else {
				state = "Not Ready"
				message = c.Message
			}
			break
		}
	}
	fmt.Fprintf(&sb, "  [lightcyan]State:[white] %s\n", state)
	if message != "" {
		fmt.Fprintf(&sb, "  [lightcyan]Message:[white] %s\n", message)
	}
	fmt.Fprintln(&sb, "")

	// DNS Records Section
	fmt.Fprintln(&sb, "[yellow::b]DNS Records[white::-]")
	for _, r := range dm.Records {
		fmt.Fprintf(&sb, "  [lightcyan]Type:[white] %s\n", r.Type)
		fmt.Fprintf(&sb, "  [lightcyan]Name:[white] %s\n", r.Name)
		fmt.Fprintf(&sb, "  [lightcyan]Data:[white] %s\n", r.RRData)
		fmt.Fprintf(&sb, "\n")
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(sb.String()).
		SetWrap(true).
		SetScrollable(true)

	textView.SetBorder(true).SetTitle(" Domain Mapping Info ")

	// Instructions
	instructions := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Press 'q' or 'esc' to close")

	// Layout
	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(instructions, 1, 0, false)

	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			closeFunc()
			return nil
		}
		return event
	})

	// Create a Grid to center the modal
	grid := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 20, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	return grid
}
