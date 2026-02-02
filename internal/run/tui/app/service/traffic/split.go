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
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "traffic-split"
)

// Modal returns a modal primitive for traffic splitting.
func Modal(app *tview.Application, service *model_service.Service, revisions []model_revision.Revision, onCompletion func(refresh bool)) tview.Primitive {
	// --- Styles ---
	fieldBackgroundColor := tcell.ColorBlack
	fieldTextColor := tcell.ColorWhite
	labelColor := tcell.ColorYellow
	buttonBgColor := tcell.ColorDarkCyan
	buttonTextColor := tcell.ColorWhite

	// --- Components ---

	// Spinner for feedback and status
	statusSpinner := spinner.New(app, 1)
	statusSpinner.SetTextAlign(tview.AlignCenter)

	// Container for Form + Status
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

	// Fields
	fields := make(map[string]*tview.InputField)
	for _, rev := range revisions {
		initialPercent := "0"
		for _, ts := range service.TrafficStatuses {
			if ts.Revision == rev.Name {
				initialPercent = strconv.Itoa(int(ts.Percent))
				break
			}
		}

		field := tview.NewInputField().
			SetLabel(rev.Name).
			SetFieldWidth(5).
			SetText(initialPercent)
		form.AddFormItem(field)
		fields[rev.Name] = field
	}

	// Buttons
	form.AddButton("Save", func() {
		var params []string
		var targets []model_traffic.TrafficTarget
		for _, rev := range revisions {
			text := fields[rev.Name].GetText()
			params = append(params, text)
			
			percent, _ := strconv.ParseInt(text, 10, 32)
			targets = append(targets, model_traffic.TrafficTarget{
				Revision: rev.Name,
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
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkGreen)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkRed)
	}

	// --- Layout ---
	container.AddItem(form, 0, 1, true)
	container.AddItem(statusSpinner, 1, 0, false)

	// Global Grid wrapper for centering
	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 4+(len(revisions)*2), 0).
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
