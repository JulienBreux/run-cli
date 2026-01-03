package loader

import (
	"fmt"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/ui/component/logo"
	"github.com/rivo/tview"
)

// New returns a new loader component.
func New(app *tview.Application) tview.Primitive {
	modal := tview.NewModal().
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	go func() {
		frames := []string{
			"Loading      Please wait",
			"Loading .    Please wait",
			"Loading ..   Please wait",
			"Loading ...  Please wait",
		}
		i := 0
		for {
			app.QueueUpdateDraw(func() {
				text := fmt.Sprintf("%s\n%s", logo.String(), frames[i])
				modal.SetText(text)
			})
			i = (i + 1) % len(frames)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return modal
}
