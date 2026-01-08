package project

import (
	"testing"

	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestProjectModal(t *testing.T) {
	app := tview.NewApplication()
	
	// Pre-populate cache
	CachedProjects = []model.Project{
		{Name: "p1"},
		{Name: "p2"},
	}
	defer func() { CachedProjects = nil }()
	
	selector := ProjectModal(app, func(p model.Project) {}, func() {})
	
	assert.NotNil(t, selector)
	assert.NotNil(t, selector.Input)
	assert.NotNil(t, selector.List)
	assert.NotNil(t, selector.Filter)
	assert.NotNil(t, selector.Submit)
	assert.NotNil(t, selector.Content)
	
	// Should satisfy Primitive interface implicitly
	var _ tview.Primitive = selector
}

func TestProjectModal_Filtering(t *testing.T) {
	app := tview.NewApplication()
	
	CachedProjects = []model.Project{
		{Name: "alpha"},
		{Name: "beta"},
		{Name: "gamma"},
	}
	defer func() { CachedProjects = nil }()
	
	selector := ProjectModal(app, func(p model.Project) {}, func() {})
	
	// Initial state: empty filter -> all items
	assert.Equal(t, 3, selector.List.GetItemCount())
	
	// Filter "a" -> alpha, beta, gamma (all contain 'a')
	selector.Filter("a")
	assert.Equal(t, 3, selector.List.GetItemCount())
	
	// Filter "al" -> alpha
	selector.Filter("al")
	assert.Equal(t, 1, selector.List.GetItemCount())
	mainText, _ := selector.List.GetItemText(0)
	assert.Equal(t, "alpha", mainText)
	
	// Filter "z" -> none
	selector.Filter("z")
	assert.Equal(t, 0, selector.List.GetItemCount())
}

func TestProjectModal_Selection(t *testing.T) {
	app := tview.NewApplication()
	
	CachedProjects = []model.Project{
		{Name: "target"},
		{Name: "other"},
	}
	defer func() { CachedProjects = nil }()
	
	var selected model.Project
	closed := false
	
	onSelect := func(p model.Project) {
		selected = p
	}
	closeModal := func() {
		closed = true
	}
	
	selector := ProjectModal(app, onSelect, closeModal)
	
	// Filter to ensure target is at index 0
	selector.Filter("target")
	assert.Equal(t, 1, selector.List.GetItemCount())
	
	// Select index 0 (which is "target")
	selector.List.SetCurrentItem(0)
	
	// Trigger submit
	selector.Submit()
	
	assert.True(t, closed)
	assert.Equal(t, "target", selected.Name)
}

func TestInputCapture(t *testing.T) {
	app := tview.NewApplication()
	closed := false
	closeModal := func() { closed = true }
	
	selector := ProjectModal(app, func(p model.Project) {}, closeModal)
	
	handler := selector.Content.GetInputCapture()
	assert.NotNil(t, handler)
	
	// Test Escape
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	ret := handler(eventEsc)
	assert.Nil(t, ret)
	assert.True(t, closed)
	
	// Test Tab Cycling
	// Initial focus is Input (set by layout order implicitly, but Application manages focus)
	// We simulate focus by setting it on Application mock? 
	// tview.Application doesn't expose GetFocus easily for verification in unit test without running.
	// However, we can verifying that SetFocus is called on the app.
	// But we passed a real app.
	
	// Since we can't easily assert "Focus changed" on `app` without internals, 
	// we will just run the handler coverage.
	
	// Simulate Input has Focus
	app.SetFocus(selector.Input)
	eventTab := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	handler(eventTab) // Should move focus to List
	assert.True(t, selector.List.HasFocus())
	
	// Simulate List has Focus
	app.SetFocus(selector.List)
	handler(eventTab) // Should move focus to BtnSelect
	// Can't check button focus easily as it is not exposed in selector.
	
	// Simulate Down arrow from Input
	app.SetFocus(selector.Input)
	eventDown := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	handler(eventDown)
	assert.True(t, selector.List.HasFocus())
}
