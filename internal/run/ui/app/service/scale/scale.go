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
	// --- Styles ---
	fieldBackgroundColor := tcell.ColorBlack
	fieldTextColor := tcell.ColorWhite
	labelColor := tcell.ColorYellow
	buttonBgColor := tcell.ColorDarkCyan
	buttonTextColor := tcell.ColorWhite

	// --- Components ---

	// Status text view for feedback
	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

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
	modeDropdown := tview.NewDropDown().
		SetLabel("Scaling mode").
		SetOptions([]string{"Automatic", "Manual"}, nil).
		SetFieldBackgroundColor(fieldBackgroundColor).
		SetListStyles(tcell.StyleDefault.Background(tcell.ColorDarkGray), tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack))

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

		// Animation control
		isSaving := true
		stopAnim := make(chan struct{})
		go func() {
			frames := []string{
				"[yellow]Saving      (Please wait)",
				"[yellow]Saving .    (Please wait)",
				"[yellow]Saving ..   (Please wait)",
				"[yellow]Saving ...  (Please wait)",
			}
			i := 0
			ticker := time.NewTicker(500 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-stopAnim:
					return
				case <-ticker.C:
					app.QueueUpdateDraw(func() {
						if isSaving {
							statusTextView.SetText(frames[i])
						}
					})
					i = (i + 1) % len(frames)
				}
			}
		}()

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()

			_, err := api_service.UpdateScaling(ctx, service.Project, service.Region, service.Name, min, max, manual)
			app.QueueUpdateDraw(func() {
				isSaving = false
				close(stopAnim)
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

	// --- Layout ---

	// Assemble Container
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusTextView, 1, 0, false)

	// Centering with Grid
	// Columns: auto, 50, auto (Centered width 50)
	// Rows: auto, 15, auto (Centered height 15)
	grid := tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 15, 0).
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// Capture escape key on the Container
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion()
			return nil
		}
		return event
	})

	return grid
}
