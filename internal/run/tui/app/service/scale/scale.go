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
func Modal(app *tview.Application, service *model_service.Service, pages *tview.Pages, onCompletion func(refresh bool)) tview.Primitive {

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

	// Container for Forms + Status
	// We use a Flex to stack ModeForm, ParamsForm, and Spinner
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).
		SetTitle(" Scale Service ").
		SetTitleAlign(tview.AlignCenter)

	// 1. Mode Form (Static)
	modeForm := tview.NewForm()
	modeForm.SetBorder(false)
	modeForm.SetLabelColor(labelColor)
	modeForm.SetFieldBackgroundColor(fieldBackgroundColor)
	modeForm.SetFieldTextColor(fieldTextColor)
	modeForm.SetItemPadding(0)

	// 2. Params Form (Dynamic)
	paramsForm := tview.NewForm()
	paramsForm.SetBorder(false)
	paramsForm.SetLabelColor(labelColor)
	paramsForm.SetFieldBackgroundColor(fieldBackgroundColor)
	paramsForm.SetFieldTextColor(fieldTextColor)
	paramsForm.SetButtonBackgroundColor(buttonBgColor)
	paramsForm.SetButtonTextColor(buttonTextColor)

	// Fields helper
	styleField := func(f *tview.InputField) {
		f.SetFieldBackgroundColor(fieldBackgroundColor)
		f.SetFieldTextColor(fieldTextColor)
	}

	// Mode Dropdown (Attached to modeForm)
	modeDropdown := dropdown.New()
	modeDropdown.SetLabel("Scaling mode")
	modeDropdown.SetOptions([]string{"Automatic", "Manual"}, nil)
	modeDropdown.SetFieldBackgroundColor(fieldBackgroundColor)
	modeDropdown.SetListStyles(
		tcell.StyleDefault.Background(tcell.ColorDarkGray),
		tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack),
	)

	// Add static items to ModeForm
	modeForm.AddFormItem(modeDropdown)

	// Form Item references for ParamsForm
	var manualInstancesField, minInstancesField, maxInstancesField *tview.InputField

	// --- Layout ---
	// ModeForm (Fixed height usually 1-2 lines) -> ParamsForm (Flex) -> Spinner
	container.AddItem(modeForm, 3, 0, true) // 3 lines default to ensure visibility
	container.AddItem(paramsForm, 0, 1, false)
	container.AddItem(statusSpinner, 1, 0, false)

	// Global Grid wrapper for centering
	grid := tview.NewGrid().
		SetColumns(0, 40, 0).
		SetRows(0, 14, 0). // Adjusted height estimate
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// --- Focus Management ---

	// Tab from ModeForm -> ParamsForm
	modeForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// If we are at the last item (which is the only item here), jump to paramsForm
			if paramsForm.GetFormItemCount() > 0 {
				app.SetFocus(paramsForm)
				return nil
			}
		}
		return event
	})

	// Backtab from ParamsForm -> ModeForm
	paramsForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyBacktab {
			// If we are at the first item, jumpt to modeForm
			index, _ := paramsForm.GetFocusedItemIndex()
			if index == 0 {
				app.SetFocus(modeForm)
				return nil
			}
		}
		return event
	})

	// Escape to Close
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion(false)
			return nil
		}
		return event
	})

	// Function to update ParamsForm based on selected mode
	updateForm := func() {
		_, mode := modeDropdown.GetCurrentOption()
		paramsForm.Clear(true) // Clear items and buttons

		if mode == "Manual" {
			val := 0
			if service.Scaling != nil {
				val = int(service.Scaling.ManualInstanceCount)
			}
			manualInstancesField = tview.NewInputField().
				SetLabel("Instances").
				SetFieldWidth(10).
				SetText(strconv.Itoa(val))
			styleField(manualInstancesField)
			paramsForm.AddFormItem(manualInstancesField)
			grid.SetRows(0, 14, 0)
		} else { // Automatic
			minVal := 0
			if service.Scaling != nil {
				minVal = int(service.Scaling.MinInstances)
			}
			minInstancesField = tview.NewInputField().
				SetLabel("Min Instances").
				SetFieldWidth(10).
				SetText(strconv.Itoa(minVal))
			styleField(minInstancesField)

			maxVal := 0
			if service.Scaling != nil {
				maxVal = int(service.Scaling.MaxInstances)
			}
			maxInstancesField = tview.NewInputField().
				SetLabel("Max Instances").
				SetFieldWidth(10).
				SetText(strconv.Itoa(maxVal))
			styleField(maxInstancesField)

			paramsForm.AddFormItem(minInstancesField)
			paramsForm.AddFormItem(maxInstancesField)
			grid.SetRows(0, 14, 0)
		}

		// Re-add Buttons
		paramsForm.AddButton("Save", func() {
			// Get values from fields safe helpers
			getText := func(f *tview.InputField) string {
				if f == nil {
					return ""
				}
				return f.GetText()
			}

			// Validate
			min, max, manual, err := validateScaleParams(
				mode,
				getText(manualInstancesField),
				getText(minInstancesField),
				getText(maxInstancesField),
			)
			if err != nil {
				statusSpinner.SetText(fmt.Sprintf("[red]%v", err))
				return
			}

			// Execute
			statusSpinner.Start("[yellow]Operation in progress...")
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
				defer cancel()
				_, err := api_service.UpdateScaling(
					ctx, service.Project,
					service.Region,
					service.Name,
					int32(min),
					int32(max),
					int32(manual),
				)
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
		paramsForm.AddButton("Cancel", func() {
			onCompletion(false)
		})

		// Style Buttons
		if paramsForm.GetButtonCount() >= 2 {
			paramsForm.GetButton(0).SetBackgroundColor(tcell.ColorDarkGreen)
			paramsForm.GetButton(1).SetBackgroundColor(tcell.ColorDarkRed)
		}
	}

	// Dropdown Selection Handler
	modeDropdown.SetSelectedFunc(func(text string, index int) {
		updateForm()
		// No need to mess with focus here; user is in ModeForm (Dropdown).
		// If they want to edit params, they press Tab.
	})

	// Initial Setup
	if service.Scaling != nil {
		if service.Scaling.ScalingMode == "MANUAL" {
			modeDropdown.SetCurrentOption(1)
			updateForm()
			manualInstancesField.SetText(strconv.Itoa(int(service.Scaling.ManualInstanceCount)))
		} else {
			modeDropdown.SetCurrentOption(0)
			updateForm()
			minInstancesField.SetText(strconv.Itoa(int(service.Scaling.MinInstances)))
			if service.Scaling.MaxInstances > 0 {
				maxInstancesField.SetText(strconv.Itoa(int(service.Scaling.MaxInstances)))
			}
		}
	} else {
		modeDropdown.SetCurrentOption(0)
		updateForm()
		minInstancesField.SetText("0")
	}

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
