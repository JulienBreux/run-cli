package spinner

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// New returns a simple "Loading..." text, right-aligned, with animation.
func New(app *tview.Application) tview.Primitive {
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignRight)

	go func() {
		frames := []string{"Loading .  ", "Loading .. ", "Loading ..."}
		i := 0
		for {
			app.QueueUpdateDraw(func() {
				textView.SetText(frames[i])
			})
			i = (i + 1) % len(frames)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return textView
}
