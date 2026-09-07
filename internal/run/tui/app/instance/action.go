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

package instance

import (
	"fmt"

	api_instance "github.com/JulienBreux/run-cli/internal/run/api/instance"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	DELETE_MODAL_PAGE_ID = "delete-instance-modal"
)

var (
	StartInstanceFunc  = api_instance.Start
	StopInstanceFunc   = api_instance.Stop
	DeleteInstanceFunc = api_instance.Delete
)

// StartAction asynchronously starts the specified instance.
func StartAction(app *tview.Application, project, region, instanceID string, onComplete func(err error)) {
	go func() {
		err := StartInstanceFunc(project, region, instanceID)
		if app != nil {
			app.QueueUpdateDraw(func() {
				if onComplete != nil {
					onComplete(err)
				}
			})
		} else if onComplete != nil {
			onComplete(err)
		}
	}()
}

// StopAction asynchronously stops the specified instance.
func StopAction(app *tview.Application, project, region, instanceID string, onComplete func(err error)) {
	go func() {
		err := StopInstanceFunc(project, region, instanceID)
		if app != nil {
			app.QueueUpdateDraw(func() {
				if onComplete != nil {
					onComplete(err)
				}
			})
		} else if onComplete != nil {
			onComplete(err)
		}
	}()
}

// TriggerDeleteConfirmation performs the async delete operation.
func TriggerDeleteConfirmation(app *tview.Application, project, region, instanceID string, onCompletion func(deleted bool, err error)) {
	go func() {
		err := DeleteInstanceFunc(project, region, instanceID)
		if app != nil {
			app.QueueUpdateDraw(func() {
				if onCompletion != nil {
					onCompletion(true, err)
				}
			})
		} else if onCompletion != nil {
			onCompletion(true, err)
		}
	}()
}

// DeleteModal creates and returns a confirmation modal for deleting an instance.
func DeleteModal(app *tview.Application, project, region, instanceID string, onCompletion func(deleted bool, err error)) *tview.Grid {
	fieldBackgroundColor := tcell.ColorBlack
	fieldTextColor := tcell.ColorWhite
	buttonBgColor := tcell.ColorDarkCyan
	buttonTextColor := tcell.ColorWhite

	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.SetBorder(true).
		SetTitle(" Delete Instance ").
		SetTitleAlign(tview.AlignCenter)

	prompt := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("\nAre you sure you want to delete instance [yellow]%s[white]\nin region [lightcyan]%s[white]?", instanceID, region))

	form := tview.NewForm()
	form.SetBorder(false)
	form.SetFieldBackgroundColor(fieldBackgroundColor)
	form.SetFieldTextColor(fieldTextColor)
	form.SetButtonBackgroundColor(buttonBgColor)
	form.SetButtonTextColor(buttonTextColor)

	form.AddButton("Delete", func() {
		TriggerDeleteConfirmation(app, project, region, instanceID, onCompletion)
	})
	form.AddButton("Cancel", func() {
		if onCompletion != nil {
			onCompletion(false, nil)
		}
	})

	if form.GetButtonCount() >= 2 {
		form.GetButton(0).SetBackgroundColor(tcell.ColorDarkRed)
		form.GetButton(1).SetBackgroundColor(tcell.ColorDarkCyan)
	}

	container.AddItem(prompt, 0, 1, false)
	container.AddItem(form, 3, 0, true)

	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 9, 0).
		AddItem(container, 1, 1, 1, 1, 0, 0, true)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if onCompletion != nil {
				onCompletion(false, nil)
			}
			return nil
		}
		return event
	})

	return grid
}
