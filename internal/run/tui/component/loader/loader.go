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

package loader

import (
	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/rivo/tview"
)

// Loader represents the loader component.
type Loader struct {
	*tview.Flex
	Spinner *spinner.Spinner
}

// New returns a new loader component.
func New(app *tview.Application) *Loader {
	// Spinner
	s := spinner.New(app, 0)
	s.SetTextAlign(tview.AlignCenter)
	s.Start("Please wait")

	// Logo
	logoView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(logo.String())

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(logoView, 6, 1, false).
		AddItem(s, 2, 0, false). // Increased height to 2 for second line
		AddItem(nil, 0, 1, false)

	return &Loader{
		Flex:    flex,
		Spinner: s,
	}
}
