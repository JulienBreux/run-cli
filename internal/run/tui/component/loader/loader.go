package loader

import (
	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/spinner"
	"github.com/rivo/tview"
)

// New returns a new loader component.
func New(app *tview.Application) tview.Primitive {
	// Spinner
	s := spinner.New(app)
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
		AddItem(s, 1, 1, false).
		AddItem(nil, 0, 1, false)

	return flex
}
