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
	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	app := tview.NewApplication()
	tbl := List(app)
	assert.NotNil(t, tbl)
	assert.Equal(t, LIST_PAGE_TITLE, tbl.Title)
}

func TestLoad(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	newInsts := []model.Instance{
		{
			Name:       "projects/p/locations/r/instances/inst-1",
			Region:     "us-central1",
			UpdateTime: time.Now(),
			TerminalCondition: &condition.Condition{
				State: "CONDITION_SUCCEEDED",
			},
			Containers: []*container.Container{
				{Name: "c1", Image: "gcr.io/demo:latest"},
			},
		},
	}

	Load(newInsts)

	assert.Equal(t, newInsts, instances)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "", listTable.Table.GetCell(1, 0).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 1).Text, "Yes")
	assert.Equal(t, "inst-1", listTable.Table.GetCell(1, 2).Text)
	assert.Equal(t, "us-central1", listTable.Table.GetCell(1, 3).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 4).Text, "Ready")
	assert.Contains(t, listTable.Table.GetCell(1, 5).Text, "c1")
}

func TestListReload(t *testing.T) {
	app := tview.NewApplication()
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init sim screen: %v", err)
	}
	app.SetScreen(simScreen)

	List(app)

	originalListFunc := ListInstancesFunc
	defer func() { ListInstancesFunc = originalListFunc }()

	expected := []model.Instance{
		{
			Name:   "projects/p/locations/r/instances/inst-reloaded",
			Region: "us-central1",
		},
	}
	ListInstancesFunc = func(projectID, region string) ([]model.Instance, error) {
		return expected, nil
	}

	ListReload(app, info.Info{Project: "p", Region: "r"}, func(err error) {
		assert.NoError(t, err)
		app.Stop()
	})

	go func() {
		time.Sleep(2 * time.Second)
		app.Stop()
	}()

	err := app.Run()
	assert.NoError(t, err)

	assert.Equal(t, expected, instances)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "inst-reloaded", listTable.Table.GetCell(1, 2).Text)
}

func TestListReload_Error(t *testing.T) {
	app := tview.NewApplication()
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init sim screen: %v", err)
	}
	app.SetScreen(simScreen)

	List(app)

	originalListFunc := ListInstancesFunc
	defer func() { ListInstancesFunc = originalListFunc }()

	ListInstancesFunc = func(projectID, region string) ([]model.Instance, error) {
		return nil, errors.New("list failed")
	}

	ListReload(app, info.Info{}, func(err error) {
		assert.Error(t, err)
		app.Stop()
	})

	go func() {
		time.Sleep(2 * time.Second)
		app.Stop()
	}()

	err := app.Run()
	assert.NoError(t, err)
}

func TestGetSelectedInstance(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	instances = []model.Instance{
		{
			Name:   "projects/p/locations/us-central1/instances/inst-1",
			Region: "us-central1",
		},
	}

	row := 1
	listTable.Table.SetCell(row, 2, tview.NewTableCell("inst-1"))
	listTable.Table.SetCell(row, 3, tview.NewTableCell("us-central1"))
	listTable.Table.Select(row, 0)

	name, region := GetSelectedInstance()
	assert.Equal(t, "inst-1", name)
	assert.Equal(t, "us-central1", region)

	full := GetSelectedInstanceFull()
	assert.NotNil(t, full)
	assert.Equal(t, "inst-1", full.ID())

	// Select Header
	listTable.Table.Select(0, 0)
	name, region = GetSelectedInstance()
	assert.Equal(t, "", name)
	assert.Equal(t, "", region)

	full = GetSelectedInstanceFull()
	assert.Nil(t, full)
}

func TestGetSelectedInstanceFull_Empty(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	instances = []model.Instance{}

	full := GetSelectedInstanceFull()
	assert.Nil(t, full)
}

func TestShortcuts(t *testing.T) {
	_ = footer.New()

	instances = []model.Instance{{Name: "projects/p/locations/r/instances/i1"}}
	assert.NotPanics(t, func() {
		Shortcuts()
	})

	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Refresh")
	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Describe")

	instances = []model.Instance{}
	Shortcuts()
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	testInsts := []model.Instance{
		{
			Name:   "projects/p/locations/r/instances/inst-1",
			Region: "us-central1",
			TerminalCondition: &condition.Condition{
				State: "CONDITION_SUCCEEDED",
			},
			Containers: []*container.Container{
				{Name: "web", Image: "nginx"},
				{Name: "sidecar", Image: "envoy"},
			},
			UpdateTime: time.Now(),
		},
		{
			Name:   "projects/p/locations/r/instances/inst-2",
			Region: "europe-west1",
			TerminalCondition: &condition.Condition{
				State: "CONDITION_FAILED",
			},
			UpdateTime: time.Now(),
		},
		{
			Name:       "projects/p/locations/r/instances/inst-3",
			Region:     "europe-west1",
			UpdateTime: time.Now(),
		},
	}

	render(testInsts)

	assert.Equal(t, 4, listTable.Table.GetRowCount())
	assert.Equal(t, "", listTable.Table.GetCell(1, 0).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 1).Text, "Yes")
	assert.Equal(t, "inst-1", listTable.Table.GetCell(1, 2).Text)
	assert.Equal(t, "us-central1", listTable.Table.GetCell(1, 3).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 4).Text, "Ready")
	assert.Contains(t, listTable.Table.GetCell(1, 5).Text, "web, sidecar")
	assert.Contains(t, listTable.Table.GetCell(2, 4).Text, "Failed")
	assert.Contains(t, listTable.Table.GetCell(3, 4).Text, "Unknown")
}

func TestRender_ProxyAndAuth(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	testInsts := []model.Instance{
		{
			Name:   "projects/p/locations/r/instances/inst-proxy-auth",
			Region: "us-central1",
			Proxy: &model_service.ProxyStatus{
				Enabled: true,
				Port:    8080,
				URL:     "http://127.0.0.1:8080",
			},
			InvokerIamDisabled: true, // Auth: No
		},
		{
			Name:               "projects/p/locations/r/instances/inst-no-proxy-auth",
			Region:             "us-central1",
			Proxy:              nil,
			InvokerIamDisabled: false, // Auth: Yes
		},
	}

	render(testInsts)

	assert.Equal(t, 3, listTable.Table.GetRowCount())
	assert.Equal(t, "[green]P", listTable.Table.GetCell(1, 0).Text)
	assert.Equal(t, "[green]No", listTable.Table.GetCell(1, 1).Text)
	assert.Equal(t, "inst-proxy-auth", listTable.Table.GetCell(1, 2).Text)
	assert.Equal(t, "us-central1", listTable.Table.GetCell(1, 3).Text)

	assert.Equal(t, "", listTable.Table.GetCell(2, 0).Text)
	assert.Equal(t, "[red]Yes", listTable.Table.GetCell(2, 1).Text)
	assert.Equal(t, "inst-no-proxy-auth", listTable.Table.GetCell(2, 2).Text)
	assert.Equal(t, "us-central1", listTable.Table.GetCell(2, 3).Text)
}

func TestGetSelectedInstanceURL(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	// No instances
	instances = []model.Instance{}
	assert.Equal(t, "", GetSelectedInstanceURL())

	// Instances with URL and Proxy
	instances = []model.Instance{
		{
			Name: "projects/p/locations/r/instances/i1",
			URLs: []string{"https://primary.run.app"},
		},
		{
			Name: "projects/p/locations/r/instances/i2",
			URLs: []string{"https://primary2.run.app"},
			Proxy: &model_service.ProxyStatus{
				Enabled: true,
				Port:    9090,
				URL:     "http://127.0.0.1:9090",
			},
		},
		{
			Name: "projects/p/locations/r/instances/i3",
		},
	}
	Load(instances)

	// Header row selected (row 0)
	listTable.Table.Select(0, 0)
	assert.Equal(t, "", GetSelectedInstanceURL())

	// Row 1 selected (URL from URLs[0])
	listTable.Table.Select(1, 0)
	assert.Equal(t, "https://primary.run.app", GetSelectedInstanceURL())

	// Row 2 selected (Proxy enabled URL)
	listTable.Table.Select(2, 0)
	assert.Equal(t, "http://127.0.0.1:9090", GetSelectedInstanceURL())

	// Row 3 selected (No URLs, no proxy)
	listTable.Table.Select(3, 0)
	assert.Equal(t, "", GetSelectedInstanceURL())
}

func TestHandleShortcuts(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testInsts := []model.Instance{
		{Name: "i1", URLs: []string{"http://test"}},
	}
	Load(testInsts)
	listTable.Table.Select(1, 0)

	// Test 'o' shortcut
	ev := tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone)
	assert.NotPanics(t, func() {
		ret := HandleShortcuts(ev)
		_ = ret
	})

	// Test unknown key
	ev2 := tcell.NewEventKey(tcell.KeyRune, 'z', tcell.ModNone)
	ret := HandleShortcuts(ev2)
	assert.Equal(t, ev2, ret)
}

func TestHandleShortcuts_ProxyOpenURL(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testInsts := []model.Instance{
		{
			Name: "i1",
			URLs: []string{"http://public"},
			Proxy: &model_service.ProxyStatus{
				Enabled: true,
				Port:    8080,
				URL:     "http://127.0.0.1:8080",
			},
		},
	}
	Load(testInsts)
	listTable.Table.Select(1, 0)

	ev := tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone)
	assert.NotPanics(t, func() {
		HandleShortcuts(ev)
	})
}

func TestHandleShortcuts_Proxy(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testInsts := []model.Instance{
		{Name: "i1", URLs: []string{"http://test"}},
	}
	Load(testInsts)
	listTable.Table.Select(1, 0)

	ev := tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone)
	assert.NotPanics(t, func() {
		HandleShortcuts(ev)
	})
}
