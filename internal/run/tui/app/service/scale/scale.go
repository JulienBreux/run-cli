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

package scale

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/JulienBreux/run-cli/pkg/dropdown"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "scale-service"
)

// Modal returns a modal primitive for scaling a service.
func Modal(app *tview.Application, service *model_service.Service, pages *tview.Pages, onCompletion func()) tview.Primitive {

	// --- Styles ---
	fieldBackgroundColor := tcell.ColorBlack
	fieldTextColor := tcell.ColorWhite
	labelColor := tcell.ColorYellow
	buttonBgColor := tcell.ColorDarkCyan
	buttonTextColor := tcell.ColorWhite

	// --- Components ---

	// Spinner for feedback and status
	statusSpinner := spinner.New(app)
	statusSpinner.SetTextAlign(tview.AlignCenter)

	// Container for Form + Status
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).
		SetTitle(" Scale Service ").
		SetTitleAlign(tview.AlignCenter)

	// Form
	form := tview.NewForm()
	form.SetBorder(false)
	form.SetLabelColor(labelColor)
	form.SetFieldBackgroundColor(fieldBackgroundColor)
	form.SetFieldTextColor(fieldTextColor)
	form.SetButtonBackgroundColor(buttonBgColor)
	form.SetButtonTextColor(buttonTextColor)

	// Fields helper
	styleField := func(f *tview.InputField) {
		f.SetFieldBackgroundColor(fieldBackgroundColor)
		f.SetFieldTextColor(fieldTextColor)
	}

	// Create form items
	var manualInstancesField, minInstancesField, maxInstancesField *tview.InputField
	modeDropdown := dropdown.New()
	modeDropdown.SetLabel("Scaling mode")
	modeDropdown.SetOptions([]string{"Automatic", "Manual"}, nil)
	modeDropdown.SetFieldBackgroundColor(fieldBackgroundColor)
	modeDropdown.SetListStyles(tcell.StyleDefault.Background(tcell.ColorDarkGray), tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack))

	manualInstancesField = tview.NewInputField().
		SetLabel("Instances").
		SetFieldWidth(10)
	styleField(manualInstancesField)

	minInstancesField = tview.NewInputField().
		SetLabel("Min Instances").
		SetFieldWidth(10)
	styleField(minInstancesField)

	maxInstancesField = tview.NewInputField().
		SetLabel("Max Instances").
		SetFieldWidth(10)
	styleField(maxInstancesField)

	// --- Layout ---

	// Assemble Container
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusSpinner, 1, 0, false)

	// Centering with Grid
	// Columns: auto, 50, auto (Centered width 50)
	// Rows: auto, 10, auto (Centered height 10)
	grid := tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 10, 0).
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// Capture escape key on the Container
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion()
			return nil
		}
		return event
	})

	// Function to update form based on selected mode
	updateForm := func() {
		_, mode := modeDropdown.GetCurrentOption()
		form.Clear(false)
		form.AddFormItem(modeDropdown)

		if mode == "Manual" {
			form.AddFormItem(manualInstancesField)
			// Rows: auto, 10, auto (Centered height 10)
			grid.SetRows(0, 10, 0)
		} else { // Automatic
			form.AddFormItem(minInstancesField)
			form.AddFormItem(maxInstancesField)
			// Rows: auto, 12, auto (Centered height 12)
			grid.SetRows(0, 12, 0)
		}
	}

	// Add buttons
	form.AddButton("Save", func() {
		// Get values from fields
		_, mode := modeDropdown.GetCurrentOption()
		min, max, manual, err := validateScaleParams(mode, manualInstancesField.GetText(), minInstancesField.GetText(), maxInstancesField.GetText())
		
		if err != nil {
			statusSpinner.SetText(fmt.Sprintf("[red]%v", err))
			return
		}

		// Start Animation
		statusSpinner.Start("[yellow]Operation in progress... (Please wait)")

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()

			_, err := api_service.UpdateScaling(ctx, service.Project, service.Region, service.Name, int32(min), int32(max), int32(manual))
			app.QueueUpdateDraw(func() {
				if err != nil {
					statusSpinner.Stop(fmt.Sprintf("[red]Error: %v", err))
				} else {
					statusSpinner.Stop("")
					onCompletion()
				}
			})
		}()
	})
	form.AddButton("Cancel", func() {
		onCompletion()
	})

	// Style Buttons (Hack: tview.Form doesn't expose buttons directly by name, so we access by index)
	// Button 0: Save (Green)
	// Button 1: Cancel (Red)
	if form.GetButtonCount() >= 2 {
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkGreen)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkRed)
	}

	// Dropdown selection handler
	modeDropdown.SetSelectedFunc(func(text string, index int) {
		updateForm()
		if text == "Manual" {
			app.SetFocus(manualInstancesField)
		} else {
			app.SetFocus(minInstancesField)
		}
	})

	// Set initial values
	if service.Scaling != nil {
		if service.Scaling.ScalingMode == "MANUAL" {
			modeDropdown.SetCurrentOption(1)
			manualInstancesField.SetText(strconv.Itoa(int(service.Scaling.ManualInstanceCount)))
		} else {
			modeDropdown.SetCurrentOption(0)
			minInstancesField.SetText(strconv.Itoa(int(service.Scaling.MinInstances)))
			if service.Scaling.MaxInstances > 0 {
				maxInstancesField.SetText(strconv.Itoa(int(service.Scaling.MaxInstances)))
			}
		}
	} else {
		// Default to Automatic
		modeDropdown.SetCurrentOption(0)
		minInstancesField.SetText("0")
	}

	updateForm() // Initial form setup

	return grid
}

func validateScaleParams(mode, manualStr, minStr, maxStr string) (min, max, manual int64, err error) {
	if mode == "Manual" {
		manual, err = strconv.ParseInt(manualStr, 10, 32)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid manual instance count")
		}
		return 0, 0, manual, nil
	}
	
	// Automatic
	min, err = strconv.ParseInt(minStr, 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid min instance count")
	}

	if maxStr != "" {
		max, err = strconv.ParseInt(maxStr, 10, 32)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid max instance count")
		}
	} else {
		max = 0
	}

	if max > 0 && min > max {
		return 0, 0, 0, fmt.Errorf("min instances cannot be greater than max instances")
	}
	
	return min, max, 0, nil
}
