package scale

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "scale"
)

// Modal returns a modal primitive for scaling a service.
func Modal(app *tview.Application, service *model_service.Service, pages *tview.Pages, onCompletion func()) tview.Primitive {
	// Main form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Service Scaling")

	// Status text view for feedback
	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Create form items
	var manualInstancesField, minInstancesField, maxInstancesField *tview.InputField
	modeDropdown := tview.NewDropDown().
		SetLabel("Scaling mode").
		SetOptions([]string{"Automatic", "Manual"}, nil)

	manualInstancesField = tview.NewInputField().
		SetLabel("Number of instances").
		SetFieldWidth(10)

	minInstancesField = tview.NewInputField().
		SetLabel("Min instances").
		SetFieldWidth(10)

	maxInstancesField = tview.NewInputField().
		SetLabel("Max instances (optional)").
		SetFieldWidth(10)

	// Function to update form based on selected mode
	updateForm := func() {
		_, mode := modeDropdown.GetCurrentOption()
		form.Clear(false)
		form.AddFormItem(modeDropdown)

		if mode == "Manual" {
			form.AddFormItem(manualInstancesField)
		} else { // Automatic
			form.AddFormItem(minInstancesField)
			form.AddFormItem(maxInstancesField)
		}
	}

	// Add buttons
	form.AddButton("Save", func() {
		// Get values from fields
		var err error
		var min, max, manual int
		_, mode := modeDropdown.GetCurrentOption()

		if mode == "Manual" {
			manual, err = strconv.Atoi(manualInstancesField.GetText())
			if err != nil {
				statusTextView.SetText("[red]Invalid manual instance count")
				return
			}
			min, max = 0, 0
		} else { // Automatic
			min, err = strconv.Atoi(minInstancesField.GetText())
			if err != nil {
				statusTextView.SetText("[red]Invalid min instance count")
				return
			}

			if maxInstancesField.GetText() != "" {
				max, err = strconv.Atoi(maxInstancesField.GetText())
				if err != nil {
					statusTextView.SetText("[red]Invalid max instance count")
					return
				}
			} else {
				max = 0
			}

			if max > 0 && min > max {
				statusTextView.SetText("[red]Min instances cannot be greater than Max instances")
				return
			}
			manual = 0
		}

		statusTextView.SetText("[yellow]Saving... (Please wait)")

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			_, err := api_service.UpdateScaling(ctx, service.Project, service.Region, service.Name, min, max, manual)
			app.QueueUpdateDraw(func() {
				if err != nil {
					statusTextView.SetText(fmt.Sprintf("[red]Error: %v", err))
				} else {
					onCompletion()
				}
			})
		}()
	})
	form.AddButton("Cancel", func() {
		onCompletion()
	})

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

	// Modal layout
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 15, 1, true).
			AddItem(statusTextView, 2, 1, false).
			AddItem(nil, 0, 1, false), 80, 1, true).
		AddItem(nil, 0, 1, false)

	// Capture escape key
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion()
			return nil
		}
		return event
	})

	return modal
}
