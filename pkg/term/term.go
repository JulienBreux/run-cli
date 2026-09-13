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

package term

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// Prefix is the standard root prefix for all CLI titles.
	Prefix = "Run"
	// Separator divides hierarchical title sections.
	Separator = " | "
)

// Output is the destination writer for terminal escape sequences.
var Output io.Writer = os.Stdout

// FormatTitle returns a formatted title string adhering to the "Run | <Part 1> | <Part 2>" pattern.
func FormatTitle(parts ...string) string {
	var filtered []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}

	if len(filtered) == 0 {
		return Prefix
	}

	if filtered[0] == Prefix {
		if len(filtered) == 1 {
			return Prefix
		}
		return strings.Join(filtered, Separator)
	}

	allParts := append([]string{Prefix}, filtered...)
	return strings.Join(allParts, Separator)
}

// SetTitle sends an ANSI OSC 0 escape sequence to set the terminal window/tab title.
func SetTitle(title string) {
	if Output != nil {
		_, _ = fmt.Fprintf(Output, "\033]0;%s\007", title)
	}
}

// ResetTitle clears the terminal title.
func ResetTitle() {
	SetTitle("")
}

// SetFormattedTitle formats the provided title parts and updates the terminal title.
func SetFormattedTitle(parts ...string) {
	SetTitle(FormatTitle(parts...))
}
