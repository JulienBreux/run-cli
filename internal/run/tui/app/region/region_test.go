package region

import (
	"testing"

	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestRegionModal_Init(t *testing.T) {
	app := tview.NewApplication()
	selector := RegionModal(app, func(s string) {}, func() {})

	assert.NotNil(t, selector)
	assert.NotNil(t, selector.Input)
	assert.NotNil(t, selector.List)
	assert.NotNil(t, selector.Filter)
	assert.NotNil(t, selector.Submit)
	
	// Should satisfy Primitive interface
	var _ tview.Primitive = selector
}

func TestRegionModal_Filtering(t *testing.T) {
	app := tview.NewApplication()
	selector := RegionModal(app, func(s string) {}, func() {})

	// Initial state: All regions + special option
	initialCount := selector.List.GetItemCount()
	assert.Greater(t, initialCount, 1)

	// Filter "us-central1"
	selector.Filter("us-central1")
	assert.Equal(t, 1, selector.List.GetItemCount())
	mainText, _ := selector.List.GetItemText(0)
	assert.Equal(t, "us-central1", mainText)

	// Filter "non-existent-region"
	selector.Filter("non-existent-region")
	assert.Equal(t, 0, selector.List.GetItemCount())
	
	// Reset
	selector.Filter("")
	assert.Equal(t, initialCount, selector.List.GetItemCount())
}

func TestRegionModal_Selection(t *testing.T) {
	app := tview.NewApplication()
	var selectedRegion string
	closed := false
	
	onSelect := func(r string) {
		selectedRegion = r
	}
	closeModal := func() {
		closed = true
	}

	selector := RegionModal(app, onSelect, closeModal)

	// Test selecting a specific region
	selector.Filter("europe-west1")
	selector.List.SetCurrentItem(0)
	selector.Submit()
	
	assert.True(t, closed)
	assert.Equal(t, "europe-west1", selectedRegion)
	
	// Reset
	closed = false
	selectedRegion = ""
	
	// Test selecting "All Regions"
	// We know "- (All Regions)" is added first in the list
	selector.Filter("") // Reset filter
	
	// Find the index of "- (All Regions)"
	idx := -1
	for i := 0; i < selector.List.GetItemCount(); i++ {
		text, _ := selector.List.GetItemText(i)
		if text == "- (All Regions)" {
			idx = i
			break
		}
	}
	assert.NotEqual(t, -1, idx, "Could not find 'All Regions' option")
	
	selector.List.SetCurrentItem(idx)
	selector.Submit()
	
	assert.True(t, closed)
	assert.Equal(t, api_region.ALL, selectedRegion)
}

func TestInputCapture(t *testing.T) {
	app := tview.NewApplication()
	closed := false
	closeModal := func() { closed = true }
	
	selector := RegionModal(app, func(s string) {}, closeModal)
	handler := selector.Content.GetInputCapture()
	
	// Test Escape
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	ret := handler(eventEsc)
	assert.Nil(t, ret)
	assert.True(t, closed)
	
	// Test Tab Cycling
	// Simulate Input has Focus
	app.SetFocus(selector.Input)
	eventTab := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	handler(eventTab) // Should move focus to List
	assert.True(t, selector.List.HasFocus())
	
	// Simulate Down arrow from Input
	app.SetFocus(selector.Input)
	eventDown := tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	handler(eventDown)
	assert.True(t, selector.List.HasFocus())
}