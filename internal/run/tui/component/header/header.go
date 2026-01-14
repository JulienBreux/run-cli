package header

import (
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/logo"
	"github.com/JulienBreux/run-cli/pkg/version"
	"github.com/rivo/tview"
)

var (
	infoView *tview.TextView
)

// New returns a TView header.
// Header composition:
// | Info | Global Shortcuts | Logo |
func New(currentInfo info.Info) *tview.Flex {
	return tview.NewFlex().
		AddItem(columnInfo(currentInfo), 50, 1, false).
		AddItem(columnShortcuts(), 0, 1, false).
		AddItem(logo.New(), 50, 1, false)
}

// UpdateInfo updates the info view.
func UpdateInfo(currentInfo info.Info) {
	infoView.Clear()

	_, _ = fmt.Fprintf(infoView, "[white]Project:        [#bd93f9]%s\n", currentInfo.Project)
	_, _ = fmt.Fprintf(infoView, "[white]Region:         [#bd93f9]%s\n", currentInfo.Region)
	_, _ = fmt.Fprintf(infoView, "[white]User:           [#bd93f9]%s\n", currentInfo.User)
	_, _ = fmt.Fprintf(infoView, "[white]Version:        [#bd93f9]%s\n", version.Version)
}

// returns the info column.
func columnInfo(currentInfo info.Info) *tview.TextView {
	infoView = tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	UpdateInfo(currentInfo)
	return infoView
}

// returns the shortcuts column.
func columnShortcuts() *tview.Flex {
	col1 := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	_, _ = fmt.Fprintf(col1, "[dodgerblue]<ctrl-p> [white]Project\n")
	_, _ = fmt.Fprintf(col1, "[dodgerblue]<ctrl-r> [white]Region\n\n")
	_, _ = fmt.Fprintf(col1, "[dodgerblue]<ctrl-z> [white]Console\n")
	_, _ = fmt.Fprintf(col1, "[dodgerblue]<ctrl-l> [white]Releases\n")

	col2 := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	_, _ = fmt.Fprintf(col2, "[dodgerblue]<ctrl-s> [white]Services\n")
	_, _ = fmt.Fprintf(col2, "[dodgerblue]<ctrl-j> [white]Jobs\n")
	_, _ = fmt.Fprintf(col2, "[dodgerblue]<ctrl-w> [white]Worker Pools\n")
	_, _ = fmt.Fprintf(col2, "[dodgerblue]<ctrl-d> [white]Domain Mappings\n")

	return tview.NewFlex().
		AddItem(col1, 20, 1, false).
		AddItem(col2, 0, 1, false)
}
