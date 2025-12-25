package header

import (
	"fmt"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/logo"
	"github.com/JulienBreux/run-cli/pkg/version"
	"github.com/rivo/tview"
)

var (
	ContextShortcutView *tview.TextView
	infoView            *tview.TextView
)

// New returns a TView header.
// Header composition:
// | Info | Shortcuts (Global and Contextual) | Logo |
func New(currentInfo info.Info) *tview.Flex {
	return tview.NewFlex().
		AddItem(columnInfo(currentInfo), 50, 1, false).
		AddItem(columnShortcuts(), 0, 1, false).
		AddItem(logo.New(), 50, 1, false)
}

// UpdateInfo updates the info view.
func UpdateInfo(currentInfo info.Info) {
	infoView.Clear()

	_, _ = fmt.Fprintf(infoView, "[white]Project:  [#bd93f9]%s\n", currentInfo.Project)
	_, _ = fmt.Fprintf(infoView, "[white]Region:   [#bd93f9]%s\n", currentInfo.Region)
	_, _ = fmt.Fprintf(infoView, "[white]User:     [#bd93f9]%s\n", currentInfo.User)
	_, _ = fmt.Fprintf(infoView, "[white]Version:  [#bd93f9]%s", version.Version)
}

// returns the info column.
func columnInfo(currentInfo info.Info) *tview.TextView {
	infoView = tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	UpdateInfo(currentInfo)
	return infoView
}

// returns the shortcuts column.
func columnShortcuts() *tview.Flex {
	return tview.NewFlex().
		AddItem(subColumnGlobalShortcuts(), 0, 1, false).
		AddItem(subColumnContextShortcuts(), 0, 1, false)
}

// returns the global shortcuts sub column.
func subColumnGlobalShortcuts() *tview.TextView {
	view := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	_, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-p> [white]Project\n")
	_, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-r> [white]Region\n\n")
	_, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-s> [white]Services\n")
	_, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-j> [white]Jobs\n")
	_, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-w> [white]Worker Pools\n")
	// _, _ = fmt.Fprintf(view, "[dodgerblue]<ctrl-i> [white]Instances\n")
	return view
}

// returns the contextual shortcuts sub column.
func subColumnContextShortcuts() *tview.TextView {
	ContextShortcutView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	return ContextShortcutView
}
