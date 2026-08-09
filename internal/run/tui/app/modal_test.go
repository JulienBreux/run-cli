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
	"os"
	"testing"

	model_project "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/credits"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/describe"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/help"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/log"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/project"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/region"
	service_scale "github.com/JulienBreux/run-cli/internal/run/tui/app/service/scale"
	workerpool_scale "github.com/JulienBreux/run-cli/internal/run/tui/app/workerpool/scale"
	"github.com/stretchr/testify/assert"
)

func TestOpenProjectModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	openProjectModal()

	assert.NotNil(t, projectModal)
	assert.Equal(t, project.MODAL_PAGE_ID, currentPageID)
	// We can't check rootPages content easily, but we can check if variable is set
}

func TestOpenRegionModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	openRegionModal()

	assert.NotNil(t, regionModal)
	assert.Equal(t, region.MODAL_PAGE_ID, currentPageID)
}

func TestModalCallbacks(t *testing.T) {
	setupTestApp()
	buildLayout()

	// Create temp home for config save
	tmpDir, _ := os.MkdirTemp("", "run-cli-app-test")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	_ = os.Setenv("HOME", tmpDir)

	// 1. Project Callback
	project.CachedProjects = []model_project.Project{{Name: "new-p"}}
	openProjectModal()
	// projectModal is a *project.ProjectSelector
	sel := projectModal.(*project.ProjectSelector)
	sel.List.SetCurrentItem(0)
	sel.Submit() // Triggers onSelect

	assert.Equal(t, "new-p", currentInfo.Project)
	assert.Equal(t, "new-p", currentConfig.Project)

	// 2. Region Callback
	openRegionModal()
	selReg := regionModal.(*region.RegionSelector)
	// Find specific region item
	selReg.Filter("us-east1")
	selReg.List.SetCurrentItem(0)
	selReg.Submit()

	assert.Equal(t, "us-east1", currentInfo.Region)
	assert.Equal(t, "us-east1", currentConfig.Region)
}

func TestOpenLogModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	openLogModal("my-service", "us-central1", "service")

	assert.Equal(t, log.MODAL_PAGE_ID, currentPageID)
}

func TestOpenDescribeModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	res := map[string]string{"foo": "bar"}
	openDescribeModal(res, "Title")

	assert.Equal(t, describe.MODAL_PAGE_ID, currentPageID)
}

func TestOpenServiceScaleModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	svc := &model_service.Service{Name: "s1"}
	openServiceScaleModal(svc)

	assert.Equal(t, service_scale.MODAL_PAGE_ID, currentPageID)
}

func TestOpenWorkerPoolScaleModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	wp := &model_workerpool.WorkerPool{DisplayName: "wp1"}
	openWorkerPoolScaleModal(wp)

	assert.Equal(t, workerpool_scale.MODAL_PAGE_ID, currentPageID)
}

func TestOpenHelpModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	openHelpModal()

	assert.Equal(t, help.MODAL_PAGE_ID, currentPageID)
}

func TestOpenCreditsModal(t *testing.T) {
	setupTestApp()
	buildLayout()

	openCreditsModal()

	assert.Equal(t, credits.MODAL_PAGE_ID, currentPageID)
}
