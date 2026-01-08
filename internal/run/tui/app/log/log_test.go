package log

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestLogModal(t *testing.T) {
	app := tview.NewApplication()
	closeModal := func() {}

	viewer := LogModal(app, "project", "filter", "My Logs", closeModal)

	assert.NotNil(t, viewer)
	assert.NotNil(t, viewer.TextView)
	assert.NotNil(t, viewer.StatusText)
	assert.NotNil(t, viewer.Content)
	
	// Should satisfy Primitive interface
	var _ tview.Primitive = viewer
}

func TestLogModal_Streaming(t *testing.T) {
	// Mock StreamLogs
	origStream := streamLogsFunc
	defer func() { streamLogsFunc = origStream }()
	
	streamLogsFunc = func(ctx context.Context, projectID, filter string, logChan chan<- string) error {
		logChan <- "Log Line 1"
		logChan <- "Log Line 2"
		// Keep channel open briefly then return? 
		// Or wait for ctx done. 
		// Real implementation blocks until done or error.
		<-ctx.Done()
		return nil
	}
	
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	go func() { _ = app.Run() }()
	defer app.Stop()
	
	viewer := LogModal(app, "p", "f", "title", func(){})
	
	// Wait for async updates
	time.Sleep(100 * time.Millisecond)
	
	text := viewer.TextView.GetText(true)
	assert.Contains(t, text, "Log Line 1")
	assert.Contains(t, text, "Log Line 2")
	
	status := viewer.StatusText.GetText(true)
	assert.Contains(t, status, "Streaming logs")
}

func TestLogModal_Error(t *testing.T) {
	// Mock StreamLogs Error
	origStream := streamLogsFunc
	defer func() { streamLogsFunc = origStream }()
	
	streamLogsFunc = func(ctx context.Context, projectID, filter string, logChan chan<- string) error {
		return errors.New("stream failed")
	}
	
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	go func() { _ = app.Run() }()
	defer app.Stop()
	
	viewer := LogModal(app, "p", "f", "title", func(){})
	
	time.Sleep(100 * time.Millisecond)
	
	text := viewer.TextView.GetText(true)
	assert.Contains(t, text, "Error streaming logs")
	assert.Contains(t, text, "stream failed")
	
	status := viewer.StatusText.GetText(true)
	assert.Equal(t, "Error", status)
}

func TestLogModal_InputCapture(t *testing.T) {
	// Mock StreamLogs to hang until cancelled
	origStream := streamLogsFunc
	defer func() { streamLogsFunc = origStream }()
	
	streamLogsFunc = func(ctx context.Context, projectID, filter string, logChan chan<- string) error {
		<-ctx.Done()
		return nil
	}
	
	app := tview.NewApplication()
	closed := false
	closeModal := func() { closed = true }
	
	viewer := LogModal(app, "p", "f", "t", closeModal)
	handler := viewer.Content.GetInputCapture()
	
	// Test Escape
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	ret := handler(eventEsc)
	
	assert.Nil(t, ret)
	assert.True(t, closed)
}