package describe

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sigs.k8s.io/yaml"
)

const MODAL_PAGE_ID = "describe"

// DescribeModal creates a modal to display resource details in YAML format.
func DescribeModal(app *tview.Application, resource any, title string, closeFunc func()) tview.Primitive {
	yamlBytes, err := yaml.Marshal(resource)
	if err != nil {
		yamlBytes = []byte("Error: Unable to marshal resource to YAML")
	}

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(string(yamlBytes)).
		SetWrap(false).
		SetScrollable(true)

	textView.SetBorder(true).SetTitle(fmt.Sprintf(" Describe: %s ", title))

	// Status/Info text
	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Press 'q' or 'esc' to close")

	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(statusText, 1, 0, false)

	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			closeFunc()
			return nil
		}
		return event
	})

	grid := tview.NewGrid().
		SetColumns(0, 120, 0).
		SetRows(0, 30, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	return grid
}
