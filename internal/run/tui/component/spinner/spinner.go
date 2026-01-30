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

package spinner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

// Spinner is a text view that displays a loading spinner.
type Spinner struct {
	*tview.TextView
	app     *tview.Application
	cancel  context.CancelFunc
	message string
	context string
	mu      sync.Mutex
	space   int
}

// New returns a new spinner component.
func New(app *tview.Application, space int) *Spinner {
	s := &Spinner{
		TextView: tview.NewTextView(),
		app:      app,
		space:    space,
	}
	s.SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true).
		SetWrap(false)
	return s
}

// SetContext sets the contextual information displayed on the second line.
func (s *Spinner) SetContext(context string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.context = context
}

// Start starts the spinner animation with the given message.
func (s *Spinner) Start(message string) {
	s.Stop("") // Stop any existing animation

	prefixSpace := strings.Repeat(" ", s.space)

	s.mu.Lock()
	s.message = message
	s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancel = cancel
	s.mu.Unlock()

	go func() {
		i := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.mu.Lock()
				msg := s.message
				ctxInfo := s.context
				s.mu.Unlock()

				s.app.QueueUpdateDraw(func() {
					text := fmt.Sprintf("%s%s %s", prefixSpace, frames[i], msg)
					if ctxInfo != "" {
						text += fmt.Sprintf("\n%s[gray]%s", prefixSpace, ctxInfo)
					}
					s.SetText(text)
				})
				i = (i + 1) % len(frames)
			}
		}
	}()
}

// Stop stops the spinner animation and sets the final message.
func (s *Spinner) Stop(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if message != "" {
		s.app.QueueUpdateDraw(func() {
			s.SetText(" " + message)
		})
	}
}
