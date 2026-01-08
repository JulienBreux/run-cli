package log

import (
	"context"
	"fmt"

	api_log "github.com/JulienBreux/run-cli/internal/run/api/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID = "modal-logs"
)

var streamLogsFunc = api_log.StreamLogs

// LogViewer represents the log viewing modal component.
type LogViewer struct {
	*tview.Grid
	Content    *tview.Flex
	TextView   *tview.TextView
	StatusText *tview.TextView
}

// LogModal returns a centered modal primitive for displaying logs
func LogModal(app *tview.Application, projectID, filter, title string, closeModal func()) *LogViewer {
	// --- Components ---

	// TextView for logs
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	textView.SetBorder(true).SetTitle(fmt.Sprintf(" Logs: %s (Streaming) ", title))

	// Status/Info text
	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Connecting to log stream...")

	// --- Logic ---

	ctx, cancel := context.WithCancel(context.Background())
	logChan := make(chan string)

	// 1. Start Streamer
	go func() {
		err := streamLogsFunc(ctx, projectID, filter, logChan)
		if err != nil {
			app.QueueUpdateDraw(func() {
				textView.SetText(fmt.Sprintf("[red]Error streaming logs: %v", err))
				statusText.SetText("Error")
			})
		}
	}()

	// 2. Start Listener
	go func() {
		for msg := range logChan {
			// Capture msg for closure
			message := msg
			app.QueueUpdateDraw(func() {
				_, _ = fmt.Fprintf(textView, "%s\n", message)
				statusText.SetText("Streaming logs... Press Esc to close.")
				textView.ScrollToEnd()
			})
		}
	}()

	// --- Layout ---

	// Main Content Flex
	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).   // Logs take all space
		AddItem(statusText, 1, 0, false) // Status line

	// --- Navigation ---
	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			cancel()       // Cancel streaming
			close(logChan) // Close channel (optional, but good practice)
			closeModal()
			return nil
		}
		return event
	})

	// --- Centering ---
	grid := tview.NewGrid().
		SetColumns(0, 120, 0).
		SetRows(0, 30, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	return &LogViewer{
		Grid:       grid,
		Content:    content,
		TextView:   textView,
		StatusText: statusText,
	}
}
