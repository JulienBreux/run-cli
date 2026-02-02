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

package logo

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

const (
	color = "[#bd93f9]"
)

// New returns a new logo component.
func New() *tview.TextView {
	logoText := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	_, _ = fmt.Fprint(logoText, String())
	return logoText
}

// String returns the logo as a string.
func String() string {
	var logo strings.Builder
	logoArt := []string{
		color + " ___ _   _ _  _ ",
		color + "| _ \\ | | | \\| |",
		color + "|   / |_| | .` |",
		color + "|_|_\\\\___/|_|\\_|",
	}

	for _, line := range logoArt {
		logo.WriteString(line)
		logo.WriteString("\n")
	}
	return logo.String()
}
