package loader

import (
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/ui/component/logo"
	"github.com/rivo/tview"
)

// New returns a new loader component.
func New() tview.Primitive {
	text := fmt.Sprintf("%s\nLoading... Please wait", logo.String())
	modal := tview.NewModal().
		SetText(text).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	return modal
}
