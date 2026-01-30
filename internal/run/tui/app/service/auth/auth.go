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

package auth

import (
	"context"
	"fmt"
	"time"

	api_service "github.com/JulienBreux/run-cli/internal/run/api/service"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/JulienBreux/run-cli/pkg/dropdown"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "auth-service"
)

// Modal returns a modal primitive for updating service authentication.
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

	// Container for Form + Status
	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).
		SetTitle(" Service Authentication ").
		SetTitleAlign(tview.AlignCenter)

	// Form
	form := tview.NewForm()
	form.SetBorder(false)
	form.SetLabelColor(labelColor)
	form.SetFieldBackgroundColor(fieldBackgroundColor)
	form.SetFieldTextColor(fieldTextColor)
	form.SetButtonBackgroundColor(buttonBgColor)
	form.SetButtonTextColor(buttonTextColor)

	// Create form items
	authDropdown := dropdown.New()
	authDropdown.SetLabel("Authentication")
	authDropdown.SetOptions([]string{"Require authentication", "Allow unauthenticated invocations"}, nil)
	authDropdown.SetFieldBackgroundColor(fieldBackgroundColor)
	authDropdown.SetListStyles(tcell.StyleDefault.Background(tcell.ColorDarkGray), tcell.StyleDefault.Background(tcell.ColorLightCyan).Foreground(tcell.ColorBlack))

	// --- Layout ---

	// Assemble Container
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusSpinner, 1, 0, false)

	// Centering with Grid
	// Columns: auto, 60, auto
	// Rows: auto, 10, auto
	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 10, 0).
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	// Capture escape key on the Container
	container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			onCompletion(false)
			return nil
		}
		return event
	})

	// Add buttons
	form.AddButton("Save", func() {
		// Get values from fields
		_, option := authDropdown.GetCurrentOption()
		allowUnauthenticated := (option == "Allow unauthenticated invocations")

		// Start Animation
		statusSpinner.Start("[yellow]Operation in progress... (Please wait)")

		// Call API
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
			defer cancel()

			_, err := api_service.UpdateAuthentication(ctx, service.Project, service.Region, service.Name, allowUnauthenticated)
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

	// Style Buttons (Hack: tview.Form doesn't expose buttons directly by name, so we access by index)
	// Button 0: Save (Green)
	// Button 1: Cancel (Red)
	if form.GetButtonCount() >= 2 {
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkGreen)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkRed)
	}

	// Set initial values
	initialOption := 0 // Require authentication
	if service.Security != nil && service.Security.InvokerIAMDisabled {
		initialOption = 1 // Allow unauthenticated invocations
	}
	authDropdown.SetCurrentOption(initialOption)

	form.AddFormItem(authDropdown)

	return grid
}
