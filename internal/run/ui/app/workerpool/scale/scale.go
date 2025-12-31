package scale

import (
	"context"
	"fmt"
	"strconv"
	"time"

	api_workerpool "github.com/JulienBreux/run-cli/internal/run/api/workerpool"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "scale-workerpool"
)

// Modal returns a modal primitive for scaling a worker pool.
func Modal(app *tview.Application, workerPool *model_workerpool.WorkerPool, pages *tview.Pages, onCompletion func()) tview.Primitive {
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
		countStr := instanceCountField.GetText()
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 0 {
			statusTextView.SetText("[red]Invalid instance count")
			return
		}

		statusTextView.SetText("[yellow]Saving... (Please wait)")

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			_, err := api_workerpool.UpdateScaling(ctx, workerPool.Project, workerPool.Region, workerPool.DisplayName, count)
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
	container.AddItem(statusTextView, 1, 0, false)

	// Centering with Grid
	// Columns: auto, 50, auto (Centered width 50)
	// Rows: auto, 10, auto (Centered height 10 - smaller than service modal)
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

	return grid
}
