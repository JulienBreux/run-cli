package dropdown

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DropDown is a wrapper around tview.DropDown to add a custom arrow indicator.
type DropDown struct {
	*tview.DropDown
}

// New creates a new DropDown.
func New() *DropDown {
	return &DropDown{
		DropDown: tview.NewDropDown(),
	}
}

// Draw draws the DropDown and the custom arrow.
func (d *DropDown) Draw(screen tcell.Screen) {
	d.DropDown.Draw(screen)

	// Custom drawing for the arrow
	x, y, width, _ := d.GetInnerRect()

	// Avoid drawing if width is too small
	if width < 2 {
		return
	}

	arrow := '▼'
	if d.IsOpen() {
		arrow = '▲'
	}

	// Calculate position: rightmost character of the first line
	targetX := x + width - 1

	// Get the existing style at the target position to preserve background color
	_, style, _ := screen.Get(targetX, y)

	// Draw the arrow with the preserved style
	screen.SetContent(targetX, y, arrow, nil, style)
}
