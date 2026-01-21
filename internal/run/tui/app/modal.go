/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"

	model_project "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	model_domainmapping "github.com/JulienBreux/run-cli/internal/run/model/domainmapping"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/credits"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/describe"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/log"
		"github.com/JulienBreux/run-cli/internal/run/tui/app/project"
		"github.com/JulienBreux/run-cli/internal/run/tui/app/region"
		service_auth "github.com/JulienBreux/run-cli/internal/run/tui/app/service/auth"
		service_scale "github.com/JulienBreux/run-cli/internal/run/tui/app/service/scale"
		workerpool_scale "github.com/JulienBreux/run-cli/internal/run/tui/app/workerpool/scale"
		"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
		"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
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
	
		footer.ContextShortcutView.Clear()
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
	
		footer.ContextShortcutView.Clear()
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
	
		footer.ContextShortcutView.Clear()
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
	
		footer.ContextShortcutView.Clear()
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
	
		footer.ContextShortcutView.Clear()
		app.SetFocus(scaleModal)
	}
	
	func openServiceAuthModal(s *model_service.Service) {
		authModal := service_auth.Modal(app, s, rootPages, func() {
			rootPages.RemovePage(service_auth.MODAL_PAGE_ID)
			switchTo(previousPageID)
		})
	
		rootPages.AddPage(service_auth.MODAL_PAGE_ID, authModal, true, true)
		previousPageID = currentPageID
		currentPageID = service_auth.MODAL_PAGE_ID
	
		footer.ContextShortcutView.Clear()
		app.SetFocus(authModal)
	}
		func openWorkerPoolScaleModal(w *model_workerpool.WorkerPool) {
		scaleModal := workerpool_scale.Modal(app, w, rootPages, func() {
			rootPages.RemovePage(workerpool_scale.MODAL_PAGE_ID)
			switchTo(previousPageID)
		})
	
		rootPages.AddPage(workerpool_scale.MODAL_PAGE_ID, scaleModal, true, true)
		previousPageID = currentPageID
		currentPageID = workerpool_scale.MODAL_PAGE_ID
	
		footer.ContextShortcutView.Clear()
		app.SetFocus(scaleModal)
	}
	
	func openCreditsModal() {
		c := credits.New(app, func() {
			rootPages.RemovePage(credits.MODAL_PAGE_ID)
			switchTo(previousPageID)
		})
	
		rootPages.AddPage(credits.MODAL_PAGE_ID, c, true, true)
		previousPageID = currentPageID
		currentPageID = credits.MODAL_PAGE_ID
	
		footer.ContextShortcutView.Clear()
		app.SetFocus(c)
		c.StartAnimation()
	}
	
	func openDomainMappingInfoModal(dm *model_domainmapping.DomainMapping) {
		modal := domainmapping.DomainMappingInfoModal(app, dm, func() {
			rootPages.RemovePage(domainmapping.MODAL_PAGE_ID)
			switchTo(previousPageID)
		})
	
		rootPages.AddPage(domainmapping.MODAL_PAGE_ID, modal, true, true)
		previousPageID = currentPageID
		currentPageID = domainmapping.MODAL_PAGE_ID
	
		footer.ContextShortcutView.Clear()
		app.SetFocus(modal)
	}

		

	