package app

import (
	"fmt"

	model_project "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/describe"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/log"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/project"
	"github.com/JulienBreux/run-cli/internal/run/ui/app/region"
	service_scale "github.com/JulienBreux/run-cli/internal/run/ui/app/service/scale"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/header"
)

func openProjectModal() {
	projectModal = project.ProjectModal(app, func(selectedProject model_project.Project) {
		currentInfo.Project = selectedProject.Name
		currentConfig.Project = selectedProject.Name
		if err := currentConfig.Save(); err != nil {
			showError(err)
			return
		}
		header.UpdateInfo(currentInfo)
	}, func() {
		rootPages.RemovePage(project.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	rootPages.AddPage(project.MODAL_PAGE_ID, projectModal, true, true)

	previousPageID = currentPageID
	currentPageID = project.MODAL_PAGE_ID

	header.ContextShortcutView.Clear()
	app.SetFocus(projectModal)
}

func openRegionModal() {
	regionModal = region.RegionModal(app, func(selectedRegion string) {
		currentInfo.Region = selectedRegion
		currentConfig.Region = selectedRegion
		if err := currentConfig.Save(); err != nil {
			showError(err)
			return
		}
		header.UpdateInfo(currentInfo)
	}, func() {
		rootPages.RemovePage(region.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	rootPages.AddPage(region.MODAL_PAGE_ID, regionModal, true, true)

	previousPageID = currentPageID
	currentPageID = region.MODAL_PAGE_ID

	header.ContextShortcutView.Clear()
	app.SetFocus(regionModal)
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
		rootPages.RemovePage(log.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	rootPages.AddPage(log.MODAL_PAGE_ID, logModal, true, true)

	previousPageID = currentPageID
	currentPageID = log.MODAL_PAGE_ID

	header.ContextShortcutView.Clear()
	app.SetFocus(logModal)
}

func openDescribeModal(resource any, title string) {
	describeModal := describe.DescribeModal(app, resource, title, func() {
		rootPages.RemovePage(describe.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	rootPages.AddPage(describe.MODAL_PAGE_ID, describeModal, true, true)

	previousPageID = currentPageID
	currentPageID = describe.MODAL_PAGE_ID

	header.ContextShortcutView.Clear()
	app.SetFocus(describeModal)
}

func openServiceScaleModal(s *model_service.Service) {
	scaleModal := service_scale.Modal(app, s, rootPages, func() {
		rootPages.RemovePage(service_scale.MODAL_PAGE_ID)
		switchTo(previousPageID)
	})

	rootPages.AddPage(service_scale.MODAL_PAGE_ID, scaleModal, true, true)
	previousPageID = currentPageID
	currentPageID = service_scale.MODAL_PAGE_ID

	header.ContextShortcutView.Clear()
	app.SetFocus(scaleModal)
}
