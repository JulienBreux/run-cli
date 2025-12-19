package logo

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

const (
	color = "[#bd93f9]"
)

// New returns a new logo component.
func New() *tview.TextView {
	logoText := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	_, _ = fmt.Fprint(logoText, String())
	return logoText
}

// String returns the logo as a string.
func String() string {
	var logo strings.Builder
	logoArt := []string{
		color + " ____  _   _ _   _  ",
		color + "|  _ \\| | | | \\ | | ",
		color + "| |_) | | | |  \\| | ",
		color + "|  _ <| |_| | . ` | ",
		color + "|_| \\_\\_____|_| \\_| ",
	}

	for _, line := range logoArt {
		logo.WriteString(line)
		logo.WriteString("\n")
	}
	logo.WriteString("\n") // Add an empty line for spacing
	logo.WriteString("[white]With[-] [red]♥[-] [white]Julien Breux[-] ")
	return logo.String()
}
