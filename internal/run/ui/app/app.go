package app

import (
	"fmt"
	"time"

	api_job "github.com/JulienBreux/run-cli/internal/run/api/job"
	"github.com/JulienBreux/run-cli/internal/run/auth"
	"github.com/JulienBreux/run-cli/internal/run/config"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/job"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/region"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/service"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/loader"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/spinner"
	"github.com/JulienBreux/run-cli/internal/run/url"
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
	currentConfig  *config.Config

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

	CONSOLE_URL = "https://console.cloud.google.com/run/overview?project=%s"
)

// Run runs the application.
func Run(cfg *config.Config) error {
	currentConfig = cfg
	app = tview.NewApplication()
	app.SetInputCapture(shortcuts)

	// Initialize default info
	currentInfo = info.Info{
		User:    "Guest",
		Project: "None",
		Region:  "all",
	}

	// Root Pages (Loader vs App)
	rootPages = tview.NewPages()
	rootPages.AddPage(LOADER_PAGE_ID, loader.New(app), true, true)

	// Start initialization in background
	go initializeApp(cfg)

	return app.SetRoot(rootPages, FULLSCREEN).
		EnableMouse(ENABLE_MOUSE).
		Run()
}

func initializeApp(cfg *config.Config) {
	// Simulate a small delay or just wait for heavy lifting
	// This helps the UI render the loader first
	time.Sleep(100 * time.Millisecond)

	// 1. Load Auth/Info (Potentially slow)
	if realInfo, err := auth.GetInfo(); err == nil {
		currentInfo.User = realInfo.User
		currentInfo.Project = realInfo.Project
		currentInfo.Region = realInfo.Region
	}

	if cfg.Region != "" {
		currentInfo.Region = cfg.Region
	}
	if cfg.Project != "" {
		currentInfo.Project = cfg.Project
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
	pages = tview.NewPages()
	// Lists
	pages.AddPage(service.LIST_PAGE_ID, service.List(app).Table, true, true)
	pages.AddPage(job.LIST_PAGE_ID, job.List(app).Table, true, true)
	pages.AddPage(workerpool.LIST_PAGE_ID, workerpool.List(app).Table, true, true)

	// Dashboards
	pages.AddPage(service.DASHBOARD_PAGE_ID, service.Dashboard(app), true, false)

	// Footer (Spinner & Error)
	errorView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	footerPages = tview.NewPages()
	footerPages.AddPage("empty", tview.NewBox(), true, true)
	footerPages.AddPage("loading", spinner.New(app), true, false)
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
	if event.Key() == tcell.KeyCtrlZ {
		u := fmt.Sprintf(CONSOLE_URL, currentInfo.Project)
		_ = url.OpenInBrowser(u)
		return nil
	}
	if event.Key() == service.LIST_PAGE_SHORTCUT {
		switchTo(service.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == job.LIST_PAGE_SHORTCUT {
		switchTo(job.LIST_PAGE_ID)
		return nil
	}
	if event.Key() == workerpool.LIST_PAGE_SHORTCUT {
		switchTo(workerpool.LIST_PAGE_ID)
		return nil
	}

	// Modals.
	if event.Key() == project.MODAL_PAGE_SHORTCUT {
		openProjectModal()
		return nil
	}
	if event.Key() == region.MODAL_PAGE_SHORTCUT {
		openRegionModal()
		return nil
	}

	if event.Key() == tcell.KeyEscape {
		if currentPageID == service.DASHBOARD_PAGE_ID {
			switchTo(service.LIST_PAGE_ID)
			return nil
		}
	}

	// Open URL for Service list
	if currentPageID == service.LIST_PAGE_ID {
		if event.Key() == tcell.KeyEnter {
			switchTo(service.DASHBOARD_PAGE_ID)
			return nil
		}
		if event.Rune() == 'r' {
			switchTo(service.LIST_PAGE_ID)
			return nil
		}
		if event.Rune() == 'l' {
			name, region := service.GetSelectedService()
			if name != "" {
				openLogModal(name, region, "service")
			}
			return nil
		}
		if event.Rune() == 'd' {
			if s := service.GetSelectedServiceFull(); s != nil {
				openDescribeModal(s, s.Name)
			}
			return nil
		}
		if event.Rune() == 's' {
			if s := service.GetSelectedServiceFull(); s != nil {
				openServiceScaleModal(s)
			}
			return nil
		}
		if result := service.HandleShortcuts(event); result == nil {
			return nil
		}
	}

	// Job List
	if currentPageID == job.LIST_PAGE_ID {
		if event.Rune() == 'r' {
			switchTo(job.LIST_PAGE_ID)
			return nil
		}
		if event.Rune() == 'l' {
			name, region := job.GetSelectedJob()
			if name != "" {
				openLogModal(name, region, "job")
			}
			return nil
		}
		if event.Rune() == 'd' {
			if j := job.GetSelectedJobFull(); j != nil {
				openDescribeModal(j, j.Name)
			}
			return nil
		}
		if event.Rune() == 'x' {
			name, region := job.GetSelectedJob()
			if name != "" {
				showLoading()
				go func() {
					_, err := api_job.Execute(currentInfo.Project, region, name)
					app.QueueUpdateDraw(func() {
						if err != nil {
							showError(err)
						} else {
							switchTo(job.LIST_PAGE_ID)
						}
					})
				}()
			}
			return nil
		}
	}

	// Worker List
	if currentPageID == workerpool.LIST_PAGE_ID {
		if event.Rune() == 'r' {
			switchTo(workerpool.LIST_PAGE_ID)
			return nil
		}
		if event.Rune() == 'd' {
			if w := workerpool.GetSelectedWorkerPoolFull(); w != nil {
				openDescribeModal(w, w.Name)
			}
			return nil
		}
		if event.Rune() == 's' {
			if w := workerpool.GetSelectedWorkerPoolFull(); w != nil {
				openWorkerPoolScaleModal(w)
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
	case service.DASHBOARD_PAGE_ID:
		if s := service.GetSelectedServiceFull(); s != nil {
			service.DashboardShortcuts()
			showLoading()
			service.DashboardReload(app, currentInfo, s, callback)
		}
	case job.LIST_PAGE_ID:
		job.Shortcuts()
		showLoading()
		job.ListReload(app, currentInfo, callback)
	case workerpool.LIST_PAGE_ID:
		workerpool.Shortcuts()
		showLoading()
		workerpool.ListReload(app, currentInfo, callback)
	}
}
