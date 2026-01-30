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

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "scale-workerpool"
)

// Modal returns a modal primitive for scaling a worker pool.
func Modal(app *tview.Application, workerPool *model_workerpool.WorkerPool, pages *tview.Pages, onCompletion func(refresh bool)) tview.Primitive {
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
		SetTitle(" Scale Worker Pool ").
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
	instanceCountField := tview.NewInputField().
		SetLabel("Instance Count").
		SetFieldWidth(10)
	styleField(instanceCountField)

	form.AddFormItem(instanceCountField)

	// Add buttons
	form.AddButton("Save", func() {
		// Get values from fields
		count, err := validateScaleParams(instanceCountField.GetText())
		if err != nil {
			statusSpinner.SetText(fmt.Sprintf("[red]%v", err))
			return
		}

		// Start Animation
		statusSpinner.Start("[yellow]Operation in progress... (Please wait)")

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			_, err := api_workerpool.UpdateScaling(ctx, workerPool.Project, workerPool.Region, workerPool.DisplayName, int32(count))
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
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkGreen)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkRed)
	}

	// Set initial values
	initialCount := 0
	if workerPool.Scaling != nil {
		initialCount = int(workerPool.Scaling.ManualInstanceCount)
	}
	instanceCountField.SetText(strconv.Itoa(initialCount))

	// --- Layout ---

	// Assemble Container
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusSpinner, 1, 0, false)

	// Centering with Grid
	// Columns: auto, 50, auto (Centered width 50)
	// Rows: auto, 8, auto (Centered height 8 - smaller than service modal)
	grid := tview.NewGrid().
		SetColumns(0, 50, 0).
		SetRows(0, 8, 0).
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// Capture escape key on the Container
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion(false)
			return nil
		}
		return event
	})

	return grid
}

func validateScaleParams(countStr string) (int64, error) {
	count, err := strconv.ParseInt(countStr, 10, 32)
	if err != nil || count < 0 {
		return 0, fmt.Errorf("invalid instance count")
	}
	return count, nil
}
