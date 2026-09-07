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

package instance

import (
	"errors"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/model/common/resources"
	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	app := tview.NewApplication()
	flex := Dashboard(app)
	assert.NotNil(t, flex)
	assert.NotNil(t, dashboardHeader)
	assert.NotNil(t, dashboardTabs)
	assert.NotNil(t, dashboardPages)
}

func TestDashboardShortcuts(t *testing.T) {
	_ = footer.New()

	assert.NotPanics(t, func() {
		DashboardShortcuts()
	})

	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Back")
}

func TestUpdateTabs(t *testing.T) {
	app := tview.NewApplication()
	_ = Dashboard(app)

	sample := &model.Instance{
		Name:        "projects/my-project/locations/us-central1/instances/inst-1",
		Region:      "us-central1",
		Project:     "my-project",
		UID:         "uid-1234",
		Generation:  1,
		CreateTime:  time.Now().Add(-1 * time.Hour),
		UpdateTime:  time.Now(),
		ServiceAccount: "sa@my-project.iam.gserviceaccount.com",
		URLs:        []string{"https://inst-1-xyz.a.run.app"},
		Labels:      map[string]string{"env": "test"},
		Annotations: map[string]string{"desc": "testing instance"},
		Containers: []*container.Container{
			{
				Name:  "web",
				Image: "gcr.io/my-project/web:v1",
				Resources: &resources.Resources{
					Limits: map[string]string{"cpu": "1", "memory": "512Mi"},
				},
			},
		},
		Conditions: []*condition.Condition{
			{
				Type:               "Ready",
				State:              "CONDITION_SUCCEEDED",
				Reason:             "ReadyReason",
				Message:            "Instance is ready",
				LastTransitionTime: time.Now(),
			},
		},
		InvokerIamDisabled: true,
	}

	dashboardInstance = sample

	assert.NotPanics(t, func() {
		updateOverviewTab()
		assert.Contains(t, overviewDetail.GetText(true), "inst-1")
		assert.Contains(t, overviewDetail.GetText(true), "my-project")
		assert.Contains(t, overviewDetail.GetText(true), "sa@my-project.iam.gserviceaccount.com")

		updateContainersTab()
		assert.Contains(t, containersDetail.GetText(true), "web")
		assert.Contains(t, containersDetail.GetText(true), "gcr.io/my-project/web:v1")
		assert.Contains(t, containersDetail.GetText(true), "512Mi")

		updateConditionsTab()
		assert.Contains(t, conditionsDetail.GetText(true), "Ready")
		assert.Contains(t, conditionsDetail.GetText(true), "Instance is ready")
		assert.Contains(t, conditionsDetail.GetText(true), "Allow unauthenticated invocations")

		updateYAMLTab()
		assert.Contains(t, yamlDetail.GetText(true), "inst-1")
	})
}

func TestDashboardReload(t *testing.T) {
	origGetFunc := GetInstanceFunc
	defer func() { GetInstanceFunc = origGetFunc }()

	called := false
	GetInstanceFunc = func(project, region, instanceID string) (*model.Instance, error) {
		called = true
		return &model.Instance{
			Name:   "projects/" + project + "/locations/" + region + "/instances/" + instanceID,
			Region: region,
		}, nil
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	_ = Dashboard(app)

	done := make(chan struct{})
	inst := &model.Instance{Name: "projects/p/locations/r/instances/i1", Region: "r"}

	go func() {
		_ = app.Run()
	}()
	defer app.Stop()

	DashboardReload(app, info.Info{Project: "p", Region: "r"}, inst, func(err error) {
		assert.NoError(t, err)
		close(done)
	})

	select {
	case <-done:
		assert.True(t, called)
		assert.NotNil(t, dashboardInstance)
		assert.Equal(t, "i1", dashboardInstance.ID())
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for DashboardReload")
	}
}

func TestDashboardReload_Error(t *testing.T) {
	origGetFunc := GetInstanceFunc
	defer func() { GetInstanceFunc = origGetFunc }()

	GetInstanceFunc = func(project, region, instanceID string) (*model.Instance, error) {
		return nil, errors.New("get instance error")
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	_ = Dashboard(app)

	done := make(chan struct{})
	inst := &model.Instance{Name: "projects/p/locations/r/instances/i1", Region: "r"}

	go func() {
		_ = app.Run()
	}()
	defer app.Stop()

	DashboardReload(app, info.Info{Project: "p", Region: "r"}, inst, func(err error) {
		assert.Error(t, err)
		close(done)
	})

	select {
	case <-done:
		assert.Nil(t, dashboardInstance)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for DashboardReload error")
	}
}

func TestDashboardInputCapture(t *testing.T) {
	app := tview.NewApplication()
	d := Dashboard(app)

	handler := d.GetInputCapture()
	assert.NotNil(t, handler)

	activeTab = 0

	// Test Tab (Next)
	eventTab := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	handler(eventTab)
	assert.Equal(t, 1, activeTab)

	// Test Right arrow (Next)
	eventRight := tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
	handler(eventRight)
	assert.Equal(t, 2, activeTab)

	// Test Backtab (Prev)
	eventBack := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
	handler(eventBack)
	assert.Equal(t, 1, activeTab)

	// Test Left arrow (Prev)
	eventLeft := tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
	handler(eventLeft)
	assert.Equal(t, 0, activeTab)

	// Test Escape Back handler
	backCalled := false
	OnDashboardBack = func() {
		backCalled = true
	}
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	handler(eventEsc)
	assert.True(t, backCalled)
}
