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

package traffic

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	model_traffic "github.com/JulienBreux/run-cli/internal/run/model/service/traffic"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/JulienBreux/run-cli/pkg/dropdown"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "traffic-split"
)

// Modal returns a modal primitive for traffic splitting using dropdowns.
func Modal(app *tview.Application, service *model_service.Service, allRevisions []model_revision.Revision, onCompletion func(refresh bool)) tview.Primitive {
	// --- Styles ---
	fieldBackgroundColor := tcell.ColorBlack
	fieldTextColor := tcell.ColorWhite
	labelColor := tcell.ColorYellow
	buttonBgColor := tcell.ColorDarkCyan
	buttonTextColor := tcell.ColorWhite

	// --- Components ---

	// Spinner
	statusSpinner := spinner.New(app, 1)
	statusSpinner.SetTextAlign(tview.AlignCenter)

	// Main Flex Container
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).
		SetTitle(" Traffic Split ").
		SetTitleAlign(tview.AlignCenter)

	// Form
	form := tview.NewForm()
	form.SetBorder(false)
	form.SetLabelColor(labelColor)
	form.SetFieldBackgroundColor(fieldBackgroundColor)
	form.SetFieldTextColor(fieldTextColor)
	form.SetButtonBackgroundColor(buttonBgColor)
	form.SetButtonTextColor(buttonTextColor)

	// Prepare Revision Options for Dropdown
	revOptions := []string{}
	for _, r := range allRevisions {
		revOptions = append(revOptions, r.Name)
	}

	// State to track rows
	type rowStruct struct {
		dropdown *dropdown.DropDown
		input    *tview.InputField
	}
	var rows []*rowStruct

	// Helper to add a row
	addRow := func(initialRev string, initialPercent string) {
		r := &rowStruct{}

		// Dropdown
		dd := dropdown.New()
		dd.SetLabel("Revision")
		dd.SetOptions(revOptions, nil)
		dd.SetFieldBackgroundColor(fieldBackgroundColor)
		dd.SetListStyles(
			tcell.StyleDefault.Background(tcell.ColorDarkGray),
			tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack),
		)
		
		// Set initial selection
		foundIndex := -1
		for i, opt := range revOptions {
			if opt == initialRev {
				foundIndex = i
				break
			}
		}
		if foundIndex >= 0 {
			dd.SetCurrentOption(foundIndex)
		} else if len(revOptions) > 0 {
			dd.SetCurrentOption(0)
		}

		// Input
		inp := tview.NewInputField().
			SetLabel("Percent").
			SetFieldWidth(5).
			SetText(initialPercent)
		inp.SetFieldBackgroundColor(fieldBackgroundColor)
		inp.SetFieldTextColor(fieldTextColor)

		r.dropdown = dd
		r.input = inp
		rows = append(rows, r)

		// Add to form
		form.AddFormItem(dd)
		form.AddFormItem(inp)
	}

	// Initial Population based on current traffic
	for _, ts := range service.TrafficStatuses {
		if ts.Percent > 0 {
			revName := ts.Revision
			if ts.Type == model_traffic.TrafficTargetAllocationTypeLatest {
				revName = service.LatestReadyRevision
			}
			if revName != "" {
				addRow(revName, strconv.Itoa(int(ts.Percent)))
			}
		}
	}
	
	// Ensure at least one row if empty (though usually traffic is 100%)
	if len(rows) == 0 {
		// Default to latest ready revision if possible, else first available
		defaultRev := service.LatestReadyRevision
		if defaultRev == "" && len(allRevisions) > 0 {
			defaultRev = allRevisions[0].Name
		}
		addRow(defaultRev, "100")
	}

	// Add Button (Dynamic Adding not easily supported by tview.Form structure in one pass, 
	// typically requires rebuilding the form or using a custom layout. 
	// For simplicity in this iteration, we will just allow editing existing splits.
	// OR we can add a "Add Split" button that rebuilds the form? 
	// tview.Form doesn't expose InsertItem easily.
	
	// Let's rely on a simpler approach: 
	// We list *all* revisions? No, that's too many.
	// We need to allow adding.
	
	// Let's try to add a "Add Revision" button *at the end*?
	// The standard Form AddButton adds to the bottom bar.
	
	form.AddButton("Add Revision", func() {
		// We can't dynamically insert form items easily into the *middle* of the rendered form list 
		// without rebuilding or hacking internals.
		// However, we can just append to the form items list.
		if len(allRevisions) > 0 {
			addRow(allRevisions[0].Name, "0")
			// We need to re-setup navigation or redraw?
			// tview handles redraw on next cycle.
		}
	})

	form.AddButton("Save", func() {
		var params []string
		var targets []model_traffic.TrafficTarget
		
		// Collect data
		// Iterate through our rows struct which holds references
		for _, r := range rows {
			_, rev := r.dropdown.GetCurrentOption()
			percentText := r.input.GetText()
			
			params = append(params, percentText)
			
			percent, _ := strconv.ParseInt(percentText, 10, 32)
			targets = append(targets, model_traffic.TrafficTarget{
				Revision: rev,
				Percent:  int32(percent),
			})
		}

		_, err := validateTrafficParams(params)
		if err != nil {
			statusSpinner.SetText(fmt.Sprintf("[red]%v", err))
			return
		}

		statusSpinner.Start("[yellow]Operation in progress...")
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()
			_, err := api_service.UpdateTraffic(ctx, service.Project, service.Region, service.Name, targets)
			app.QueueUpdateDraw(func() {
				if err != nil {
					statusSpinner.Stop(fmt.Sprintf("[red]Error: %v", err))
				} else {
					statusSpinner.Stop("")
					onCompletion(true)
				}
			})
		}()
	})

	form.AddButton("Cancel", func() {
		onCompletion(false)
	})

	// Style Buttons
	if form.GetButtonCount() >= 2 {
		// 0: Add, 1: Save, 2: Cancel
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkBlue)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkGreen)
		form.GetButton(2).SetBackgroundColor(tcell.ColorDarkRed)
	}

	// --- Layout ---
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusSpinner, 1, 0, false)

	// Global Grid wrapper for centering
	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 20, 0). // Fixed height for now, scrolling handled by Form if needed? Form scrolls? 
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// Escape to Close
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion(false)
			return nil
		}
		return event
	})

	return grid
}

func validateTrafficParams(params []string) (int64, error) {
	var sum int64
	for _, p := range params {
		val, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid percentage: %s", p)
		}
		if val < 0 || val > 100 {
			return 0, fmt.Errorf("percentage must be between 0 and 100: %s", p)
		}
		sum += val
	}
	if sum != 100 {
		return sum, fmt.Errorf("total percentage must be 100, current: %d", sum)
	}
	return sum, nil
}