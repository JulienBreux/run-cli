package region

import (
	"strings"

	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	MODAL_PAGE_ID       = "modal-regions"
	MODAL_PAGE_SHORTCUT = tcell.KeyCtrlR
)

// RegionSelector represents the region selection modal component.
type RegionSelector struct {
	*tview.Grid
	Content *tview.Flex
	Input   *tview.InputField
	List    *tview.List
	Filter  func(string)
	Submit  func()
}

// RegionModal returns a centered modal primitive with search and list
func RegionModal(app *tview.Application, onSelect func(region string), closeModal func()) *RegionSelector {
	// --- Data ---
	regions := append([]string{"- (All Regions)"}, api_region.List()...)
	var filteredRegions []string

	// --- Components ---

	// Input
	input := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30).
		SetLabelColor(tcell.ColorYellow)

	// List
	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(tcell.ColorDarkBlue)

	list.SetBorder(true).SetTitle(" Results ")

	// Buttons
	btnSelect := tview.NewButton("Select").SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkGreen))
	btnCancel := tview.NewButton("Cancel").SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed))

	// --- Logic ---

	populateList := func(filter string) {
		list.Clear()
		filteredRegions = nil
		filter = strings.ToLower(filter)
		for _, r := range regions {
			if strings.Contains(strings.ToLower(r), filter) {
				filteredRegions = append(filteredRegions, r)
				list.AddItem(r, "", 0, nil)
			}
		}
	}

	// Init List
	populateList("")

	// Events
	input.SetChangedFunc(populateList)

	submit := func() {
		idx := list.GetCurrentItem()
		if idx != -1 && idx < len(filteredRegions) {
			selection := filteredRegions[idx]
			// Handle "All Regions" special case
			if selection == "- (All Regions)" {
				onSelect(api_region.ALL)
			} else {
				onSelect(selection)
			}
			closeModal()
		}
	}

	btnSelect.SetSelectedFunc(submit)
	btnCancel.SetSelectedFunc(closeModal)

	// Allow selecting items directly from the list
	list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		submit()
	})

	// --- Layout (The Box) ---

	// 1. Flex for Buttons
	buttons := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(btnSelect, 12, 1, false).
		AddItem(nil, 2, 0, false). // Space between buttons
		AddItem(btnCancel, 12, 1, false).
		AddItem(nil, 0, 1, false)

	// 2. Main Content Flex (Vertical)
	content := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(input, 1, 0, true).   // Search Bar
		AddItem(nil, 1, 0, false).    // Padding
		AddItem(list, 0, 1, false).   // List (Takes remaining space in the box)
		AddItem(nil, 1, 0, false).    // Padding
		AddItem(buttons, 1, 0, false) // Buttons

	content.SetBorder(true).
		SetTitle(" Select Region ").
		SetTitleAlign(tview.AlignCenter)

	// --- Navigation (Tab Cycling) ---
	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			closeModal()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			switch {
			case input.HasFocus():
				app.SetFocus(list)
			case list.HasFocus():
				app.SetFocus(btnSelect)
			case btnSelect.HasFocus():
				app.SetFocus(btnCancel)
			case btnCancel.HasFocus():
				app.SetFocus(input)
			}
			return nil
		}
		// Convenience: Down arrow from Input goes to List
		if input.HasFocus() && event.Key() == tcell.KeyDown {
			app.SetFocus(list)
			return nil
		}
		return event
	})

	// --- Centering (The Grid) ---
	// We use a Grid to create a perfectly centered float.
	// Columns: 0 (flexible), 60 (fixed width for modal), 0 (flexible)
	// Rows:    0 (flexible), 20 (fixed height for modal), 0 (flexible)
	grid := tview.NewGrid().
		SetColumns(0, 60, 0).
		SetRows(0, 20, 0).
		AddItem(content, 1, 1, 1, 1, 0, 0, true)

	return &RegionSelector{
		Grid:    grid,
		Content: content,
		Input:   input,
		List:    list,
		Filter:  populateList,
		Submit:  submit,
	}
}
