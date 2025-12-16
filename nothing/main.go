package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// 1. Background Layer
	background := tview.NewTextView().
		SetText("Background. Press 'O' to open the Project Modal.").
		SetTextAlign(tview.AlignCenter)
	pages.AddPage("background", background, true, true)

	// 2. Global Input Capture (Open Modal)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'o' || event.Rune() == 'O' {
			// Check if modal is already open
			if pages.HasPage("modal") {
				return event
			}

			// Create the modal primitive
			modal := ProjectModal(app, pages, func(selected string) {
				background.SetText(fmt.Sprintf("You selected: %s", selected))
			})

			// Add modal layer
			pages.AddPage("modal", modal, true, true)

			// FORCE FOCUS: We must explicitly tell the app to look at the new modal
			app.SetFocus(modal)
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// ProjectModal creates the popup
func ProjectModal(app *tview.Application, pages *tview.Pages, onSelect func(string)) tview.Primitive {
	// --- Data ---
	projects := make([]string, 35)
	for i := 0; i < 35; i++ {
		projects[i] = fmt.Sprintf("Project Alpha %02d", i+1)
	}

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
		filter = strings.ToLower(filter)
		for _, p := range projects {
			if strings.Contains(strings.ToLower(p), filter) {
				list.AddItem(p, "", 0, nil)
			}
		}
	}

	// Init List
	populateList("")

	// Events
	input.SetChangedFunc(populateList)

	closeModal := func() {
		pages.RemovePage("modal")
	}

	submit := func() {
		if list.GetCurrentItem() != -1 {
			text, _ := list.GetItemText(list.GetCurrentItem())
			onSelect(text)
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
		SetTitle(" Select Project ").
		SetTitleAlign(tview.AlignCenter)

	// --- Navigation (Tab Cycling) ---
	content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	return grid
}
