package app

import (
	"fmt"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/auth"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_project "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/job"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/log"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/region"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/service"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/worker"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/loader"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/spinner"
	"github.com/JulienBreux/run-cli/internal/run/ui/header"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app       *tview.Application
	rootPages *tview.Pages // Root pages to hold Layout and Loader
	pages     *tview.Pages // Content pages (Service, Job, etc.)

	previousPageID string
	currentPageID  string
	currentInfo    info.Info

	projectModal tview.Primitive
	regionModal  tview.Primitive

	footerPages *tview.Pages
	errorView   *tview.TextView
)

const (
	FULLSCREEN   = true
	ENABLE_MOUSE = true

	ESCAPE_SHORTCUT = tcell.KeyEscape
	LOADER_PAGE_ID  = "loader"
	LAYOUT_PAGE_ID  = "layout"
)

// Run runs the application.
func Run() error {
	app = tview.NewApplication()
	app.SetInputCapture(shortcuts)

	// Initialize default info
	currentInfo = info.Info{
		User:    "Guest",
		Project: "None",
		Version: "dev",
		Region:  "-",
	}

	// Root Pages (Loader vs App)
	rootPages = tview.NewPages()
	rootPages.AddPage(LOADER_PAGE_ID, loader.New(), true, true)

	// Start initialization in background
	go initializeApp()

	return app.SetRoot(rootPages, FULLSCREEN).
		EnableMouse(ENABLE_MOUSE).
		Run()
}

func initializeApp() {
	// Simulate a small delay or just wait for heavy lifting
	// This helps the UI render the loader first
	time.Sleep(100 * time.Millisecond)

	// 1. Load Auth/Info (Potentially slow)
	if realInfo, err := auth.GetInfo(); err == nil {
		currentInfo.User = realInfo.User
		currentInfo.Project = realInfo.Project
		currentInfo.Region = realInfo.Region
	}

	app.QueueUpdateDraw(func() {
		// 2. Build the main layout
		mainLayout := buildLayout()
		rootPages.AddPage(LAYOUT_PAGE_ID, mainLayout, true, false)

		// 3. Switch to main layout
		rootPages.SwitchToPage(LAYOUT_PAGE_ID)

		// 4. Trigger initial data load
		switchTo(service.LIST_PAGE_ID)
	})
}

// buildLayout constructs the main application UI
func buildLayout() *tview.Flex {
	// Modals
	projectModal = project.ProjectModal(app, func(selectedProject model_project.Project) {
		currentInfo.Project = selectedProject.Name
		header.UpdateInfo(currentInfo)
		switchTo(previousPageID)
	}, func() {
		switchTo(service.LIST_PAGE_ID)
	})

	regionModal = region.RegionModal(app, func(selectedRegion string) {
		currentInfo.Region = selectedRegion
		header.UpdateInfo(currentInfo)
		switchTo(previousPageID)
	}, func() {
		switchTo(service.LIST_PAGE_ID)
	})

	pages = tview.NewPages()
	// Lists
	pages.AddPage(service.LIST_PAGE_ID, service.List(app).Table, true, true)
	pages.AddPage(job.LIST_PAGE_ID, job.List(app).Table, true, true)
	pages.AddPage(worker.LIST_PAGE_ID, worker.List(app).Table, true, true)

	// Modals to Pages
	pages.AddPage(project.MODAL_PAGE_ID, projectModal, true, true)
	pages.AddPage(region.MODAL_PAGE_ID, regionModal, true, true)

	// Footer (Spinner & Error)
	errorView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	footerPages = tview.NewPages()
	footerPages.AddPage("empty", tview.NewBox(), true, true)
	footerPages.AddPage("loading", spinner.New(), true, false)
	footerPages.AddPage("error", errorView, true, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header.New(currentInfo), 7, 1, false).
		AddItem(pages, 0, 1, true).
		AddItem(footerPages, 1, 1, false)

	return layout
}

// shortcuts captures all key events.
func shortcuts(event *tcell.EventKey) *tcell.EventKey {
	// Disable shortcuts if we are still loading
	frontPage, _ := rootPages.GetFrontPage()
	if frontPage == LOADER_PAGE_ID {
		return nil
	}

	// Navigation.
	if event.Key() == service.LIST_PAGE_SHORTCUT {
		switchTo(service.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == job.LIST_PAGE_SHORTCUT {
		switchTo(job.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == worker.LIST_PAGE_SHORTCUT {
		switchTo(worker.LIST_PAGE_ID)
		return nil
	}

	// Modals.
	if event.Key() == project.MODAL_PAGE_SHORTCUT {
		switchTo(project.MODAL_PAGE_ID)
		return nil
	}
	if event.Key() == region.MODAL_PAGE_SHORTCUT {
		switchTo(region.MODAL_PAGE_ID)
		return nil
	}

	// Open URL for Service list
	if currentPageID == service.LIST_PAGE_ID {
		if event.Rune() == 'l' {
			name, region := service.GetSelectedService()
			if name != "" {
				openLogModal(name, region, "service")
			}
			return nil
		}
		if result := service.HandleShortcuts(event); result == nil {
			return nil
		}
	}

	// Job List
	if currentPageID == job.LIST_PAGE_ID {
		if event.Rune() == 'l' {
			name, region := job.GetSelectedJob()
			if name != "" {
				openLogModal(name, region, "job")
			}
			return nil
		}
	}

	return event
}

func showLoading() {
	footerPages.SwitchToPage("loading")
}

func hideLoading() {
	footerPages.SwitchToPage("empty")
}

func showError(err error) {
	errorView.SetText(fmt.Sprintf("[red]%s", err.Error()))
	footerPages.SwitchToPage("error")
}

func switchTo(pageID string) {
	previousPageID = currentPageID
	currentPageID = pageID
	pages.SwitchToPage(pageID)

	callback := func(err error) {
		if err != nil {
			showError(err)
		} else {
			hideLoading()
		}
	}

	switch pageID {
	case service.LIST_PAGE_ID:
		service.Shortcuts()
		showLoading()
		service.ListReload(app, currentInfo, callback)
	case job.LIST_PAGE_ID:
		job.Shortcuts()
		showLoading()
		job.ListReload(app, currentInfo, callback)
	case worker.LIST_PAGE_ID:
		worker.Shortcuts()
		showLoading()
		worker.ListReload(app, currentInfo, callback)
	case project.MODAL_PAGE_ID:
		header.ContextShortcutView.Clear()
		app.SetFocus(projectModal)
	case region.MODAL_PAGE_ID:
		header.ContextShortcutView.Clear()
		app.SetFocus(regionModal)
	}
}

func openLogModal(name, region, logType string) {
	var filter string
	switch logType {
	case "service":
		filter = fmt.Sprintf(`resource.type="cloud_run_revision" resource.labels.service_name="%s" resource.labels.location="%s"`, name, region)
	case "job":
		filter = fmt.Sprintf(`resource.type="cloud_run_job" resource.labels.job_name="%s" resource.labels.location="%s"`, name, region)
	}

	logModal := log.LogModal(app, currentInfo.Project, filter, name, func() {
		pages.RemovePage(log.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	pages.AddPage(log.MODAL_PAGE_ID, logModal, true, true)

	previousPageID = currentPageID
	currentPageID = log.MODAL_PAGE_ID
	pages.SwitchToPage(log.MODAL_PAGE_ID)

	header.ContextShortcutView.Clear()
	app.SetFocus(logModal)
}
